// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Deployers} from "v4-core/test/utils/Deployers.sol";
import {IHooks} from "v4-core/src/interfaces/IHooks.sol";
import {PoolId, PoolIdLibrary} from "v4-core/src/types/PoolId.sol";
import {PoolKey} from "v4-core/src/types/PoolKey.sol";
import {Currency} from "v4-core/src/types/Currency.sol";
import {Vm} from "forge-std/Vm.sol";
import {Create2} from "@openzeppelin/contracts/utils/Create2.sol";

import {FeeEscrow} from "../src/FeeEscrow.sol";
import {LaunchHookV2} from "../src/LaunchHookV2.sol";
import {LaunchpadRegistry} from "../src/LaunchpadRegistry.sol";
import {LaunchToken} from "../src/LaunchToken.sol";
import {MultiTenantLaunchpadFactory} from "../src/MultiTenantLaunchpadFactory.sol";

contract MultiTenantLaunchpadFactoryTest is Deployers {
    using PoolIdLibrary for PoolKey;

    LaunchpadRegistry internal registry;
    LaunchHookV2 internal hook;
    FeeEscrow internal escrow;
    MultiTenantLaunchpadFactory internal factory;

    address internal owner = makeAddr("launchpadOwner");
    address internal treasury = makeAddr("launchpadTreasury");
    address internal creator = makeAddr("tokenCreator");
    address internal laasTreasury = makeAddr("laasTreasury");
    bytes32 internal launchpadId;

    function setUp() public {
        deployFreshManagerAndRouters();
        deployMintAndApprove2Currencies();
        registry = new LaunchpadRegistry();

        LaunchHookV2 implementation = new LaunchHookV2(manager, address(this));
        uint160 flags = (1 << 13) | (1 << 11) | (1 << 9) | (1 << 7) | (1 << 6) | (1 << 3) | (1 << 2);
        address hookAddress = address(flags);
        vm.etch(hookAddress, address(implementation).code);
        hook = LaunchHookV2(hookAddress);
        escrow = new FeeEscrow(manager, hookAddress);

        factory =
            new MultiTenantLaunchpadFactory(registry, manager, hook, Currency.unwrap(currency1), laasTreasury, 0, 60);
        hook.setFactoryOnce(address(factory), address(escrow));

        vm.prank(owner);
        launchpadId = registry.createLaunchpad("ai-pad", treasury);
    }

    function testRejectsUnknownAndInactiveLaunchpad() public {
        MultiTenantLaunchpadFactory.LaunchParams memory params = _params();
        params.launchpadId = keccak256("unknown");
        vm.expectRevert(MultiTenantLaunchpadFactory.InvalidLaunchpad.selector);
        vm.prank(creator);
        factory.launch(params);

        vm.prank(owner);
        registry.setActive(launchpadId, false);
        params.launchpadId = launchpadId;
        vm.expectRevert(MultiTenantLaunchpadFactory.InvalidLaunchpad.selector);
        vm.prank(creator);
        factory.launch(params);
    }

    function testRejectsMinimumStartTickFrame() public {
        vm.expectRevert(MultiTenantLaunchpadFactory.InvalidConfig.selector);
        new MultiTenantLaunchpadFactory(
            registry, manager, hook, Currency.unwrap(currency1), laasTreasury, type(int24).min, 60
        );
    }

    function testRejectsStaleConfigAndExpiredDeadline() public {
        MultiTenantLaunchpadFactory.LaunchParams memory params = _params();
        params.expectedConfigVersion = 2;
        vm.expectRevert(abi.encodeWithSelector(MultiTenantLaunchpadFactory.StaleConfig.selector, uint64(2), uint64(1)));
        vm.prank(creator);
        factory.launch(params);

        params.expectedConfigVersion = 1;
        params.deadline = uint64(block.timestamp - 1);
        vm.expectRevert(abi.encodeWithSelector(MultiTenantLaunchpadFactory.LaunchExpired.selector, params.deadline));
        vm.prank(creator);
        factory.launch(params);
    }

    function testLaunchesFixedSupplyTokenForAnotherWallet() public {
        vm.recordLogs();
        vm.prank(creator);
        (address token, bytes32 poolId) = factory.launch(_params());

        assertEq(LaunchToken(token).totalSupply(), 1_000_000_000 ether);
        assertEq(factory.launchpadOf(token), launchpadId);
        assertEq(factory.poolOf(token), poolId);
        assertTrue(hook.seeded(PoolId.wrap(poolId)));
        assertEq(LaunchToken(token).balanceOf(address(hook)), 0);

        (
            bool initialized,
            ,
            address frozenCreator,
            address protocolTreasury,
            address frozenLaasTreasury,
            uint16 protocolFeeBps,
            uint16 laasFeeBps,
            uint16 creatorBps,
            uint16 platformBps,
            uint16 referrerBps,
            uint16 antiSnipeStartTotalBps,
            uint32 antiSnipeWindowSeconds,
            uint48 launchTime
        ) = hook.poolConfig(poolId);
        assertTrue(initialized);
        assertEq(frozenCreator, creator);
        assertEq(protocolTreasury, treasury);
        assertEq(frozenLaasTreasury, laasTreasury);
        assertEq(protocolFeeBps, 100);
        assertEq(laasFeeBps, 50);
        assertEq(creatorBps, 5000);
        assertEq(platformBps, 3000);
        assertEq(referrerBps, 2000);
        assertEq(antiSnipeStartTotalBps, 9900);
        assertEq(antiSnipeWindowSeconds, 16);
        assertEq(launchTime, block.timestamp);

        Vm.Log[] memory logs = vm.getRecordedLogs();
        bytes32 launchedSignature = keccak256("Launched(address,bytes32,bytes32,address,address,uint256,int24)");
        bool found;
        for (uint256 i; i < logs.length; ++i) {
            if (logs[i].emitter == address(factory) && logs[i].topics[0] == launchedSignature) {
                assertEq(logs[i].topics[1], bytes32(uint256(uint160(token))));
                assertEq(logs[i].topics[2], poolId);
                assertEq(logs[i].topics[3], launchpadId);
                found = true;
            }
        }
        assertTrue(found);
    }

    function testRegistryTreasuryChangeOnlyAffectsFuturePools() public {
        vm.prank(creator);
        (, bytes32 firstPoolId) = factory.launch(_params());

        address nextTreasury = makeAddr("nextTreasury");
        vm.prank(owner);
        registry.setTreasury(launchpadId, nextTreasury);

        MultiTenantLaunchpadFactory.LaunchParams memory next = _params();
        next.name = "Beta";
        next.symbol = "BETA";
        next.salt = keccak256("beta-salt");
        vm.prank(creator);
        (, bytes32 secondPoolId) = factory.launch(next);

        (,,, address firstTreasury,,,,,,,,,) = hook.poolConfig(firstPoolId);
        (,,, address secondTreasury,,,,,,,,,) = hook.poolConfig(secondPoolId);
        assertEq(firstTreasury, treasury);
        assertEq(secondTreasury, nextTreasury);
    }

    function testRejectsSaltReplay() public {
        MultiTenantLaunchpadFactory.LaunchParams memory params = _params();
        vm.prank(creator);
        factory.launch(params);

        bytes32 derivedSalt = keccak256(abi.encode(creator, params.salt));
        vm.expectRevert(abi.encodeWithSelector(MultiTenantLaunchpadFactory.LaunchSaltUsed.selector, derivedSalt));
        vm.prank(creator);
        factory.launch(params);
    }

    function testCreate2PredictionAndBothCurrencyOrderings() public {
        MultiTenantLaunchpadFactory.LaunchParams memory token0Params = _paramsForOrdering(true, "Token Zero", "ZERO");
        address predictedToken0 = _predictedToken(token0Params, creator);
        vm.prank(creator);
        (address token0, bytes32 pool0) = factory.launch(token0Params);

        assertEq(token0, predictedToken0);
        assertLt(uint160(token0), uint160(Currency.unwrap(currency1)));
        assertEq(pool0, _expectedPoolId(token0, true));

        MultiTenantLaunchpadFactory.LaunchParams memory token1Params = _paramsForOrdering(false, "Token One", "ONE");
        address predictedToken1 = _predictedToken(token1Params, creator);
        vm.prank(creator);
        (address token1, bytes32 pool1) = factory.launch(token1Params);

        assertEq(token1, predictedToken1);
        assertGt(uint160(token1), uint160(Currency.unwrap(currency1)));
        assertEq(pool1, _expectedPoolId(token1, false));
    }

    function testCreate2SaltIsBoundToCaller() public {
        MultiTenantLaunchpadFactory.LaunchParams memory params = _params();
        address anotherCreator = makeAddr("anotherCreator");
        address creatorPrediction = _predictedToken(params, creator);
        address anotherPrediction = _predictedToken(params, anotherCreator);
        assertNotEq(creatorPrediction, anotherPrediction);

        vm.prank(creator);
        (address creatorToken,) = factory.launch(params);
        vm.prank(anotherCreator);
        (address anotherToken,) = factory.launch(params);

        assertEq(creatorToken, creatorPrediction);
        assertEq(anotherToken, anotherPrediction);
    }

    function _params() internal view returns (MultiTenantLaunchpadFactory.LaunchParams memory) {
        return MultiTenantLaunchpadFactory.LaunchParams({
            launchpadId: launchpadId,
            name: "Alpha",
            symbol: "ALPHA",
            contractURI: "ipfs://alpha",
            salt: keccak256("alpha-salt"),
            quote: Currency.unwrap(currency1),
            expectedConfigVersion: 1,
            deadline: uint64(block.timestamp + 1 hours)
        });
    }

    function _paramsForOrdering(bool tokenIsCurrency0, string memory name, string memory symbol)
        internal
        view
        returns (MultiTenantLaunchpadFactory.LaunchParams memory params)
    {
        params = _params();
        params.name = name;
        params.symbol = symbol;
        params.contractURI = string.concat("ipfs://", symbol);

        for (uint256 i; i < 512; ++i) {
            params.salt = bytes32(i);
            bool predictedIsCurrency0 = uint160(_predictedToken(params, creator)) < uint160(params.quote);
            if (predictedIsCurrency0 == tokenIsCurrency0) return params;
        }
        revert("ordering salt not found");
    }

    function _predictedToken(MultiTenantLaunchpadFactory.LaunchParams memory params, address caller)
        internal
        view
        returns (address)
    {
        bytes32 derivedSalt = keccak256(abi.encode(caller, params.salt));
        bytes memory initCode = abi.encodePacked(
            type(LaunchToken).creationCode,
            abi.encode(params.name, params.symbol, params.contractURI, factory.LAUNCH_SUPPLY(), address(hook))
        );
        return Create2.computeAddress(derivedSalt, keccak256(initCode), address(factory));
    }

    function _expectedPoolId(address token, bool tokenIsCurrency0) internal view returns (bytes32) {
        (Currency currency0_, Currency currency1_) =
            tokenIsCurrency0 ? (Currency.wrap(token), currency1) : (currency1, Currency.wrap(token));
        PoolKey memory key = PoolKey({
            currency0: currency0_,
            currency1: currency1_,
            fee: 0,
            tickSpacing: factory.tickSpacing(),
            hooks: IHooks(address(hook))
        });
        return PoolId.unwrap(key.toId());
    }
}
