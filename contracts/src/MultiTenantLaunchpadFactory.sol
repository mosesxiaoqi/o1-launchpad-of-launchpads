// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Create2} from "@openzeppelin/contracts/utils/Create2.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {IPoolManager} from "v4-core/src/interfaces/IPoolManager.sol";
import {IHooks} from "v4-core/src/interfaces/IHooks.sol";
import {PoolKey} from "v4-core/src/types/PoolKey.sol";
import {PoolId, PoolIdLibrary} from "v4-core/src/types/PoolId.sol";
import {Currency} from "v4-core/src/types/Currency.sol";
import {TickMath} from "v4-core/src/libraries/TickMath.sol";
import {LiquidityAmounts} from "v4-core/test/utils/LiquidityAmounts.sol";

import {LaunchHookV2} from "./LaunchHookV2.sol";
import {LaunchpadRegistry} from "./LaunchpadRegistry.sol";
import {LaunchToken} from "./LaunchToken.sol";

contract MultiTenantLaunchpadFactory is ReentrancyGuard {
    using PoolIdLibrary for PoolKey;

    uint256 public constant LAUNCH_SUPPLY = 1_000_000_000 ether;
    uint64 public constant configVersion = 1;

    LaunchpadRegistry public immutable registry;
    IPoolManager public immutable poolManager;
    LaunchHookV2 public immutable hook;
    address public immutable quote;
    address public immutable laasTreasury;
    int24 public immutable startTickToken0Frame;
    int24 public immutable tickSpacing;

    mapping(bytes32 salt => bool used) public usedLaunchSalts;
    mapping(address token => bytes32 launchpadId) public launchpadOf;
    mapping(address token => bytes32 poolId) public poolOf;

    struct LaunchParams {
        bytes32 launchpadId;
        string name;
        string symbol;
        string contractURI;
        bytes32 salt;
        address quote;
        uint64 expectedConfigVersion;
        uint64 deadline;
    }

    error EmptyName();
    error InvalidConfig();
    error InvalidLaunchpad();
    error LaunchExpired(uint64 deadline);
    error LaunchSaltUsed(bytes32 salt);
    error StaleConfig(uint64 expectedVersion, uint64 actualVersion);

    event Launched(
        address indexed token,
        bytes32 indexed poolId,
        bytes32 indexed launchpadId,
        address creator,
        address quote,
        uint256 supply,
        int24 tickSpacing
    );

    constructor(
        LaunchpadRegistry registry_,
        IPoolManager poolManager_,
        LaunchHookV2 hook_,
        address quote_,
        address laasTreasury_,
        int24 startTickToken0Frame_,
        int24 tickSpacing_
    ) {
        if (
            address(registry_) == address(0) || address(poolManager_) == address(0) || address(hook_) == address(0)
                || laasTreasury_ == address(0) || startTickToken0Frame_ == type(int24).min || tickSpacing_ <= 0
                || tickSpacing_ > 16_384
        ) revert InvalidConfig();
        registry = registry_;
        poolManager = poolManager_;
        hook = hook_;
        quote = quote_;
        laasTreasury = laasTreasury_;
        startTickToken0Frame = startTickToken0Frame_;
        tickSpacing = tickSpacing_;
    }

    function launch(LaunchParams calldata params) external nonReentrant returns (address token, bytes32 poolId) {
        if (params.expectedConfigVersion != configVersion) {
            revert StaleConfig(params.expectedConfigVersion, configVersion);
        }
        if (block.timestamp > params.deadline) revert LaunchExpired(params.deadline);
        if (bytes(params.name).length == 0 || bytes(params.symbol).length == 0) revert EmptyName();
        if (params.quote != quote) revert InvalidConfig();

        (address launchpadOwner, address protocolTreasury, bool active,) = registry.getLaunchpad(params.launchpadId);
        if (launchpadOwner == address(0) || protocolTreasury == address(0) || !active) revert InvalidLaunchpad();

        bytes32 derivedSalt = keccak256(abi.encode(msg.sender, params.salt));
        if (usedLaunchSalts[derivedSalt]) revert LaunchSaltUsed(derivedSalt);
        usedLaunchSalts[derivedSalt] = true;

        bytes memory initCode = abi.encodePacked(
            type(LaunchToken).creationCode,
            abi.encode(params.name, params.symbol, params.contractURI, LAUNCH_SUPPLY, address(hook))
        );
        address predicted = Create2.computeAddress(derivedSalt, keccak256(initCode), address(this));
        bool tokenIsCurrency0 = uint160(predicted) < uint160(quote);

        int24 spacing = tickSpacing;
        int24 frame = tokenIsCurrency0 ? startTickToken0Frame : -startTickToken0Frame;
        int24 startTick = (frame / spacing) * spacing;
        int24 maxUsable = (TickMath.MAX_TICK / spacing) * spacing;
        int24 minUsable = (TickMath.MIN_TICK / spacing) * spacing;
        int24 tickLower = tokenIsCurrency0 ? startTick : minUsable;
        int24 tickUpper = tokenIsCurrency0 ? maxUsable : startTick;
        int24 initTick = tokenIsCurrency0 ? startTick - spacing : startTick;
        if (
            tickLower >= tickUpper || initTick < TickMath.MIN_TICK || initTick > TickMath.MAX_TICK
                || tickLower % spacing != 0 || tickUpper % spacing != 0
        ) revert InvalidConfig();

        uint160 sqrtLower = TickMath.getSqrtPriceAtTick(tickLower);
        uint160 sqrtUpper = TickMath.getSqrtPriceAtTick(tickUpper);
        uint128 liquidity = tokenIsCurrency0
            ? LiquidityAmounts.getLiquidityForAmount0(sqrtLower, sqrtUpper, LAUNCH_SUPPLY)
            : LiquidityAmounts.getLiquidityForAmount1(sqrtLower, sqrtUpper, LAUNCH_SUPPLY);
        if (liquidity == 0 || liquidity >= uint128(type(int128).max)) revert InvalidConfig();

        token = address(
            new LaunchToken{salt: derivedSalt}(
                params.name, params.symbol, params.contractURI, LAUNCH_SUPPLY, address(hook)
            )
        );
        if (token != predicted) revert InvalidConfig();

        (Currency currency0, Currency currency1) = tokenIsCurrency0
            ? (Currency.wrap(token), Currency.wrap(quote))
            : (Currency.wrap(quote), Currency.wrap(token));
        PoolKey memory key = PoolKey({
            currency0: currency0,
            currency1: currency1,
            fee: 0,
            tickSpacing: spacing,
            hooks: IHooks(address(hook))
        });

        poolManager.initialize(key, TickMath.getSqrtPriceAtTick(initTick));
        hook.registerPool(
            key,
            LaunchHookV2.PoolConfig({
                initialized: false,
                tokenIsCurrency0: tokenIsCurrency0,
                creator: msg.sender,
                protocolTreasury: protocolTreasury,
                laasTreasury: laasTreasury,
                protocolFeeBps: 100,
                laasFeeBps: 50,
                creatorBps: 5000,
                platformBps: 3000,
                referrerBps: 2000,
                antiSnipeStartTotalBps: 9900,
                antiSnipeWindowSeconds: 16,
                launchTime: uint48(block.timestamp)
            })
        );
        LaunchHookV2.SeedPosition[] memory positions = new LaunchHookV2.SeedPosition[](1);
        positions[0] = LaunchHookV2.SeedPosition({tickLower: tickLower, tickUpper: tickUpper, liquidity: liquidity});
        hook.seedLiquidity(key, positions, tokenIsCurrency0);

        poolId = PoolId.unwrap(key.toId());
        launchpadOf[token] = params.launchpadId;
        poolOf[token] = poolId;
        emit Launched(token, poolId, params.launchpadId, msg.sender, quote, LAUNCH_SUPPLY, spacing);
    }
}
