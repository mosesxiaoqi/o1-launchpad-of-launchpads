// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IHooks} from "v4-core/src/interfaces/IHooks.sol";
import {IPoolManager} from "v4-core/src/interfaces/IPoolManager.sol";
import {Hooks} from "v4-core/src/libraries/Hooks.sol";
import {SafeCast} from "v4-core/src/libraries/SafeCast.sol";
import {PoolKey} from "v4-core/src/types/PoolKey.sol";
import {PoolId, PoolIdLibrary} from "v4-core/src/types/PoolId.sol";
import {Currency} from "v4-core/src/types/Currency.sol";
import {BalanceDelta, BalanceDeltaLibrary} from "v4-core/src/types/BalanceDelta.sol";
import {BeforeSwapDelta, BeforeSwapDeltaLibrary, toBeforeSwapDelta} from "v4-core/src/types/BeforeSwapDelta.sol";
import {ModifyLiquidityParams, SwapParams} from "v4-core/src/types/PoolOperation.sol";
import {IERC20Minimal} from "v4-core/src/interfaces/external/IERC20Minimal.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";

import {FeeEscrow} from "./FeeEscrow.sol";
import {BaseHook} from "./base/BaseHook.sol";

contract LaunchHookV2 is BaseHook {
    using PoolIdLibrary for PoolKey;
    using BalanceDeltaLibrary for BalanceDelta;
    using SafeCast for uint256;

    uint16 public constant PROTOCOL_FEE_BPS = 100;
    uint16 public constant LAAS_FEE_BPS = 50;
    uint16 public constant NORMAL_TOTAL_FEE_BPS = PROTOCOL_FEE_BPS + LAAS_FEE_BPS;
    uint16 public constant CREATOR_SHARE_BPS = 5000;
    uint16 public constant PLATFORM_SHARE_BPS = 3000;
    uint16 public constant REFERRER_SHARE_BPS = 2000;
    uint16 public constant MAX_TOTAL_FEE_BPS = 9900;
    uint32 public constant ANTI_SNIPE_WINDOW_SECONDS = 16;
    uint256 internal constant BPS = 10_000;

    struct PoolConfig {
        bool initialized;
        bool tokenIsCurrency0;
        address creator;
        address protocolTreasury;
        address laasTreasury;
        uint16 protocolFeeBps;
        uint16 laasFeeBps;
        uint16 creatorBps;
        uint16 platformBps;
        uint16 referrerBps;
        uint16 antiSnipeStartTotalBps;
        uint32 antiSnipeWindowSeconds;
        uint48 launchTime;
    }

    struct FeeSplit {
        uint256 creator;
        uint256 protocol;
        uint256 referrer;
        uint256 laas;
        uint256 surcharge;
        uint256 total;
    }

    struct SeedPosition {
        int24 tickLower;
        int24 tickUpper;
        uint128 liquidity;
    }

    error InvalidFeeConfig();
    error AlreadyWired();
    error AlreadySeeded();
    error BadWiring();
    error LiquidityLocked();
    error NotDeployer();
    error NotFactory();
    error NotSingleSided();
    error NotWired();
    error ExactOutputDisabledDuringAntiSnipe();
    error PartialFillUnsupported();
    error UnexpectedFeeCurrency();
    error PoolAlreadyRegistered();
    error UnknownPool();
    error ZeroAddress();

    address public immutable deployer;
    address public factory;
    address public feeEscrow;

    mapping(bytes32 poolId => PoolConfig config) public poolConfig;
    mapping(PoolId poolId => bool value) public seeded;

    event Trade(
        bytes32 indexed poolId,
        address indexed executor,
        address indexed referrer,
        address feeCurrency,
        uint256 totalFee,
        bytes32 comment
    );

    constructor(IPoolManager poolManager_, address deployer_) BaseHook(poolManager_) {
        if (address(poolManager_) == address(0) || deployer_ == address(0)) revert ZeroAddress();
        deployer = deployer_;
    }

    function getHookPermissions() public pure override returns (Hooks.Permissions memory) {
        return Hooks.Permissions({
            beforeInitialize: true,
            afterInitialize: false,
            beforeAddLiquidity: true,
            afterAddLiquidity: false,
            beforeRemoveLiquidity: true,
            afterRemoveLiquidity: false,
            beforeSwap: true,
            afterSwap: true,
            beforeDonate: false,
            afterDonate: false,
            beforeSwapReturnDelta: true,
            afterSwapReturnDelta: true,
            afterAddLiquidityReturnDelta: false,
            afterRemoveLiquidityReturnDelta: false
        });
    }

    function setFactoryOnce(address factory_, address feeEscrow_) external {
        if (msg.sender != deployer) revert NotDeployer();
        if (factory != address(0)) revert AlreadyWired();
        if (factory_ == address(0) || feeEscrow_ == address(0)) revert ZeroAddress();
        FeeEscrow escrow = FeeEscrow(feeEscrow_);
        if (escrow.hook() != address(this) || address(escrow.poolManager()) != address(poolManager)) {
            revert BadWiring();
        }
        factory = factory_;
        feeEscrow = feeEscrow_;
    }

    function registerPool(PoolKey calldata key, PoolConfig calldata config) external {
        if (msg.sender != factory) revert NotFactory();
        PoolConfig memory frozen = config;
        frozen.initialized = true;
        _registerPoolConfig(PoolId.unwrap(key.toId()), frozen);
    }

    function seedLiquidity(PoolKey calldata key, SeedPosition[] calldata positions, bool tokenIsCurrency0) external {
        if (msg.sender != factory) revert NotFactory();
        PoolId poolId = key.toId();
        PoolConfig storage config = poolConfig[PoolId.unwrap(poolId)];
        if (!config.initialized) revert UnknownPool();
        if (seeded[poolId]) revert AlreadySeeded();
        if (config.tokenIsCurrency0 != tokenIsCurrency0 || positions.length == 0) revert NotSingleSided();

        seeded[poolId] = true;
        poolManager.unlock(abi.encode(key, positions, tokenIsCurrency0));
    }

    function unlockCallback(bytes calldata data) external returns (bytes memory) {
        if (msg.sender != address(poolManager)) revert NotPoolManager();
        (PoolKey memory key, SeedPosition[] memory positions, bool tokenIsCurrency0) =
            abi.decode(data, (PoolKey, SeedPosition[], bool));

        Currency tokenCurrency = tokenIsCurrency0 ? key.currency0 : key.currency1;
        uint256 tokenOwed;
        for (uint256 i; i < positions.length; ++i) {
            (BalanceDelta delta,) = poolManager.modifyLiquidity(
                key,
                ModifyLiquidityParams({
                    tickLower: positions[i].tickLower,
                    tickUpper: positions[i].tickUpper,
                    liquidityDelta: int256(uint256(positions[i].liquidity)),
                    salt: bytes32(0)
                }),
                ""
            );

            int128 quoteDelta = tokenIsCurrency0 ? delta.amount1() : delta.amount0();
            int128 tokenDelta = tokenIsCurrency0 ? delta.amount0() : delta.amount1();
            if (quoteDelta != 0 || tokenDelta >= 0) revert NotSingleSided();
            tokenOwed += _magnitude(tokenDelta);
        }

        if (Currency.unwrap(tokenCurrency) == address(0)) {
            poolManager.settle{value: tokenOwed}();
        } else {
            poolManager.sync(tokenCurrency);
            IERC20Minimal(Currency.unwrap(tokenCurrency)).transfer(address(poolManager), tokenOwed);
            poolManager.settle();
        }
        return "";
    }

    function currentTotalFeeBps(bytes32 poolId, uint256 timestamp) public view returns (uint256) {
        PoolConfig memory config = poolConfig[poolId];
        if (!config.initialized) revert UnknownPool();
        return _currentTotalFeeBps(config, timestamp);
    }

    function _currentTotalFeeBps(PoolConfig memory config, uint256 timestamp) private pure returns (uint256) {
        if (timestamp <= config.launchTime) return config.antiSnipeStartTotalBps;

        uint256 elapsed = timestamp - config.launchTime;
        if (elapsed >= config.antiSnipeWindowSeconds) return NORMAL_TOTAL_FEE_BPS;

        uint256 surchargeAtLaunch = uint256(config.antiSnipeStartTotalBps) - NORMAL_TOTAL_FEE_BPS;
        uint256 surcharge =
            surchargeAtLaunch * (config.antiSnipeWindowSeconds - elapsed) / config.antiSnipeWindowSeconds;
        return NORMAL_TOTAL_FEE_BPS + surcharge;
    }

    function previewFeeSplit(bytes32 poolId, uint256 amount, address trader, address referrer, uint256 timestamp)
        external
        view
        returns (FeeSplit memory split)
    {
        PoolConfig memory config = poolConfig[poolId];
        if (!config.initialized) revert UnknownPool();

        uint256 protocolFee = amount * config.protocolFeeBps / BPS;
        split.laas = amount * config.laasFeeBps / BPS;
        split.total = amount * _currentTotalFeeBps(config, timestamp) / BPS;
        split.creator = protocolFee * config.creatorBps / BPS;

        bool validReferrer = _isValidReferrer(config, trader, referrer);
        split.referrer = validReferrer ? protocolFee * config.referrerBps / BPS : 0;
        split.protocol = protocolFee - split.creator - split.referrer;
        split.surcharge = split.total - protocolFee - split.laas;
    }

    function _registerPoolConfig(bytes32 poolId, PoolConfig memory config) internal {
        if (poolConfig[poolId].initialized) revert PoolAlreadyRegistered();
        if (
            !config.initialized || config.creator == address(0) || config.protocolTreasury == address(0)
                || config.laasTreasury == address(0) || config.protocolFeeBps != PROTOCOL_FEE_BPS
                || config.laasFeeBps != LAAS_FEE_BPS || config.creatorBps != CREATOR_SHARE_BPS
                || config.platformBps != PLATFORM_SHARE_BPS || config.referrerBps != REFERRER_SHARE_BPS
                || config.antiSnipeStartTotalBps != MAX_TOTAL_FEE_BPS
                || config.antiSnipeWindowSeconds != ANTI_SNIPE_WINDOW_SECONDS
        ) revert InvalidFeeConfig();

        poolConfig[poolId] = config;
    }

    function _beforeInitialize(address sender, PoolKey calldata, uint160) internal view override returns (bytes4) {
        if (factory == address(0)) revert NotWired();
        if (sender != factory) revert NotFactory();
        return IHooks.beforeInitialize.selector;
    }

    function _beforeAddLiquidity(address, PoolKey calldata, ModifyLiquidityParams calldata, bytes calldata)
        internal
        pure
        override
        returns (bytes4)
    {
        revert LiquidityLocked();
    }

    function _beforeRemoveLiquidity(address, PoolKey calldata, ModifyLiquidityParams calldata, bytes calldata)
        internal
        pure
        override
        returns (bytes4)
    {
        revert LiquidityLocked();
    }

    function _beforeSwap(address, PoolKey calldata key, SwapParams calldata params, bytes calldata)
        internal
        override
        returns (bytes4, BeforeSwapDelta, uint24)
    {
        bytes32 poolId = PoolId.unwrap(key.toId());
        PoolConfig memory config = poolConfig[poolId];
        if (!config.initialized) return (IHooks.beforeSwap.selector, BeforeSwapDeltaLibrary.ZERO_DELTA, 0);

        uint256 totalFeeBps = _currentTotalFeeBps(config, block.timestamp);
        if (params.amountSpecified > 0 && totalFeeBps > NORMAL_TOTAL_FEE_BPS) {
            revert ExactOutputDisabledDuringAntiSnipe();
        }

        (Currency quoteCurrency, bool quoteIsSpecified) = _quoteCurrencyAndSpecified(key, params, config);
        if (!quoteIsSpecified) return (IHooks.beforeSwap.selector, BeforeSwapDeltaLibrary.ZERO_DELTA, 0);

        bool exactOutput = params.amountSpecified > 0;
        uint256 totalFee = _feeAmount(_specifiedMagnitude(params.amountSpecified), totalFeeBps, exactOutput);
        if (totalFee == 0) return (IHooks.beforeSwap.selector, BeforeSwapDeltaLibrary.ZERO_DELTA, 0);

        poolManager.mint(feeEscrow, quoteCurrency.toId(), totalFee);
        return (IHooks.beforeSwap.selector, toBeforeSwapDelta(totalFee.toInt128(), 0), 0);
    }

    function _afterSwap(
        address sender,
        PoolKey calldata key,
        SwapParams calldata params,
        BalanceDelta delta,
        bytes calldata hookData
    ) internal override returns (bytes4, int128) {
        bytes32 poolId = PoolId.unwrap(key.toId());
        PoolConfig memory config = poolConfig[poolId];
        if (!config.initialized) return (IHooks.afterSwap.selector, int128(0));

        uint256 totalFeeBps = _currentTotalFeeBps(config, block.timestamp);
        if (params.amountSpecified > 0 && totalFeeBps > NORMAL_TOTAL_FEE_BPS) {
            revert ExactOutputDisabledDuringAntiSnipe();
        }

        (Currency quoteCurrency, bool quoteIsSpecified) = _quoteCurrencyAndSpecified(key, params, config);
        bool exactOutput = params.amountSpecified > 0;
        if (quoteIsSpecified) {
            uint256 specifiedMagnitude = _specifiedMagnitude(params.amountSpecified);
            uint256 specifiedTotalFee = _feeAmount(specifiedMagnitude, totalFeeBps, exactOutput);
            if (specifiedTotalFee != 0) {
                _requireFullSpecifiedFill(params, delta, specifiedTotalFee);
                _distribute(
                    config, quoteCurrency, specifiedMagnitude, specifiedTotalFee, exactOutput, hookData, poolId, sender
                );
            }
            return (IHooks.afterSwap.selector, int128(0));
        }

        (Currency feeCurrency, uint256 unspecifiedMagnitude) = _unspecified(key, params, delta);
        if (Currency.unwrap(feeCurrency) != Currency.unwrap(quoteCurrency)) revert UnexpectedFeeCurrency();
        if (unspecifiedMagnitude == 0) return (IHooks.afterSwap.selector, int128(0));

        uint256 unspecifiedTotalFee = _feeAmount(unspecifiedMagnitude, totalFeeBps, exactOutput);
        if (unspecifiedTotalFee == 0) return (IHooks.afterSwap.selector, int128(0));
        poolManager.mint(feeEscrow, feeCurrency.toId(), unspecifiedTotalFee);
        _distribute(
            config, feeCurrency, unspecifiedMagnitude, unspecifiedTotalFee, exactOutput, hookData, poolId, sender
        );
        return (IHooks.afterSwap.selector, unspecifiedTotalFee.toInt128());
    }

    function _distribute(
        PoolConfig memory config,
        Currency feeCurrency,
        uint256 magnitude,
        uint256 totalFee,
        bool exactOutput,
        bytes calldata hookData,
        bytes32 poolId,
        address sender
    ) private {
        uint256 protocolFee;
        uint256 laasFee;
        if (exactOutput) {
            uint256 normalFee = _feeAmount(magnitude, NORMAL_TOTAL_FEE_BPS, true);
            protocolFee = normalFee * PROTOCOL_FEE_BPS / NORMAL_TOTAL_FEE_BPS;
            laasFee = normalFee - protocolFee;
        } else {
            protocolFee = magnitude * PROTOCOL_FEE_BPS / BPS;
            laasFee = magnitude * LAAS_FEE_BPS / BPS;
        }

        uint256 creatorShare = protocolFee * config.creatorBps / BPS;
        (address referrer, bytes32 comment) = _parseHookData(hookData);
        bool validReferrer = _isValidReferrer(config, sender, referrer);
        uint256 referrerShare = validReferrer ? protocolFee * config.referrerBps / BPS : 0;
        uint256 protocolShare = totalFee - creatorShare - referrerShare - laasFee;
        address currency = Currency.unwrap(feeCurrency);

        FeeEscrow escrow = FeeEscrow(payable(feeEscrow));
        if (creatorShare != 0) escrow.credit(config.creator, currency, creatorShare);
        if (protocolShare != 0) escrow.credit(config.protocolTreasury, currency, protocolShare);
        if (referrerShare != 0) escrow.credit(referrer, currency, referrerShare);
        if (laasFee != 0) escrow.credit(config.laasTreasury, currency, laasFee);
        emit Trade(poolId, sender, validReferrer ? referrer : address(0), currency, totalFee, comment);
    }

    function _isValidReferrer(PoolConfig memory config, address trader, address referrer) private pure returns (bool) {
        return referrer != address(0) && referrer != trader && referrer != config.creator
            && referrer != config.protocolTreasury && referrer != config.laasTreasury;
    }

    function _feeAmount(uint256 magnitude, uint256 feeBps, bool exactOutput) private pure returns (uint256) {
        if (feeBps == 0) return 0;
        if (!exactOutput) return Math.mulDiv(magnitude, feeBps, BPS);
        return Math.mulDiv(magnitude, feeBps, BPS - feeBps, Math.Rounding.Ceil);
    }

    function _quoteCurrencyAndSpecified(PoolKey calldata key, SwapParams calldata params, PoolConfig memory config)
        private
        pure
        returns (Currency quoteCurrency, bool quoteIsSpecified)
    {
        quoteCurrency = config.tokenIsCurrency0 ? key.currency1 : key.currency0;
        bool exactInput = params.amountSpecified < 0;
        bool specifiedIsCurrency0 = exactInput ? params.zeroForOne : !params.zeroForOne;
        Currency specifiedCurrency = specifiedIsCurrency0 ? key.currency0 : key.currency1;
        quoteIsSpecified = Currency.unwrap(specifiedCurrency) == Currency.unwrap(quoteCurrency);
    }

    function _specifiedMagnitude(int256 amountSpecified) private pure returns (uint256) {
        if (amountSpecified < 0) return uint256(-(amountSpecified + 1)) + 1;
        return uint256(amountSpecified);
    }

    function _magnitude(int128 amount) internal pure returns (uint256) {
        if (amount < 0) return uint256(uint128(-(amount + 1))) + 1;
        return uint256(uint128(amount));
    }

    function _requireFullSpecifiedFill(SwapParams calldata params, BalanceDelta delta, uint256 specifiedFee)
        private
        pure
    {
        bool exactInput = params.amountSpecified < 0;
        bool specifiedIsCurrency0 = exactInput ? params.zeroForOne : !params.zeroForOne;
        int128 actualSpecified = specifiedIsCurrency0 ? delta.amount0() : delta.amount1();
        if (int256(actualSpecified) != params.amountSpecified + int256(specifiedFee)) {
            revert PartialFillUnsupported();
        }
    }

    function _unspecified(PoolKey calldata key, SwapParams calldata params, BalanceDelta delta)
        private
        pure
        returns (Currency currency, uint256 magnitude)
    {
        bool exactInput = params.amountSpecified < 0;
        bool unspecifiedIsCurrency1 = exactInput ? params.zeroForOne : !params.zeroForOne;
        int128 amount;
        if (unspecifiedIsCurrency1) {
            currency = key.currency1;
            amount = delta.amount1();
        } else {
            currency = key.currency0;
            amount = delta.amount0();
        }
        magnitude = _magnitude(amount);
    }

    function _parseHookData(bytes calldata hookData) internal pure returns (address referrer, bytes32 comment) {
        if (hookData.length >= 32) referrer = address(uint160(uint256(bytes32(hookData[:32]))));
        if (hookData.length >= 64) comment = bytes32(hookData[32:64]);
    }
}
