// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {MockERC20} from "solmate/src/test/utils/mocks/MockERC20.sol";
import {Deployers} from "v4-core/test/utils/Deployers.sol";
import {IHooks} from "v4-core/src/interfaces/IHooks.sol";
import {IPoolManager} from "v4-core/src/interfaces/IPoolManager.sol";
import {PoolKey} from "v4-core/src/types/PoolKey.sol";
import {PoolId, PoolIdLibrary} from "v4-core/src/types/PoolId.sol";
import {Currency} from "v4-core/src/types/Currency.sol";
import {StateLibrary} from "v4-core/src/libraries/StateLibrary.sol";
import {Hooks} from "v4-core/src/libraries/Hooks.sol";

import {FeeEscrow} from "../src/FeeEscrow.sol";
import {LaunchHookV2} from "../src/LaunchHookV2.sol";

contract LaunchHookV2IntegrationTest is Deployers {
    using StateLibrary for IPoolManager;
    using PoolIdLibrary for PoolKey;

    LaunchHookV2 internal hook;
    FeeEscrow internal escrow;
    PoolKey internal launchKey;
    PoolId internal launchPoolId;

    function setUp() public {
        deployFreshManagerAndRouters();
        deployMintAndApprove2Currencies();

        LaunchHookV2 implementation = new LaunchHookV2(manager, address(this));
        uint160 flags = (1 << 13) | (1 << 11) | (1 << 9) | (1 << 7) | (1 << 6) | (1 << 3) | (1 << 2);
        address hookAddress = address(flags);
        vm.etch(hookAddress, address(implementation).code);
        hook = LaunchHookV2(hookAddress);
        escrow = new FeeEscrow(manager, hookAddress);
        hook.setFactoryOnce(address(this), address(escrow));

        launchKey =
            PoolKey({currency0: currency0, currency1: currency1, fee: 0, tickSpacing: 60, hooks: IHooks(hookAddress)});
        launchPoolId = launchKey.toId();
        manager.initialize(launchKey, SQRT_PRICE_1_1);
        hook.registerPool(launchKey, _defaultConfig());
        MockERC20(Currency.unwrap(currency0)).transfer(hookAddress, 1_000_000 ether);
    }

    function testSeedsSingleSidedPositionOnce() public {
        LaunchHookV2.SeedPosition[] memory positions = new LaunchHookV2.SeedPosition[](1);
        positions[0] = LaunchHookV2.SeedPosition({tickLower: 60, tickUpper: 120, liquidity: 1e18});

        hook.seedLiquidity(launchKey, positions, true);

        (uint128 liquidity,,) = manager.getPositionInfo(launchPoolId, address(hook), 60, 120, bytes32(0));
        assertEq(liquidity, 1e18);
        assertTrue(hook.seeded(launchPoolId));

        vm.expectRevert(LaunchHookV2.AlreadySeeded.selector);
        hook.seedLiquidity(launchKey, positions, true);
    }

    function testHookAddressMatchesDeclaredPermissions() public view {
        Hooks.validateHookPermissions(IHooks(address(hook)), hook.getHookPermissions());
    }

    function testExactInputBuyCreditsProtocolAndLaaSShares() public {
        _seed();
        vm.warp(block.timestamp + 16);
        address creator = makeAddr("creator");
        address protocolTreasury = makeAddr("protocolTreasury");
        address laasTreasury = makeAddr("laasTreasury");
        address referrer = makeAddr("referrer");

        swap(launchKey, false, -1000, abi.encode(referrer, bytes32(0)));

        address quote = Currency.unwrap(currency1);
        assertEq(escrow.owed(creator, quote), 5);
        assertEq(escrow.owed(protocolTreasury, quote), 3);
        assertEq(escrow.owed(referrer, quote), 2);
        assertEq(escrow.owed(laasTreasury, quote), 5);
        assertEq(manager.balanceOf(address(escrow), uint160(quote)), 15);
        assertEq(MockERC20(quote).balanceOf(address(escrow)), 0);
        assertEq(escrow.totalOwed(quote), 15);

        escrow.claim(creator, quote);
        assertEq(MockERC20(quote).balanceOf(creator), 5);
        assertEq(manager.balanceOf(address(escrow), uint160(quote)), 10);
        assertEq(escrow.totalOwed(quote), 10);
    }

    function testExactInputSellChargesQuoteOutput() public {
        _seed();
        vm.warp(block.timestamp + 16);
        swap(launchKey, false, -10_000, "");
        address quote = Currency.unwrap(currency1);
        uint256 owedAfterBuy = escrow.totalOwed(quote);

        swap(launchKey, true, -1000, "");

        assertGt(escrow.totalOwed(quote), owedAfterBuy);
        assertEq(manager.balanceOf(address(escrow), uint160(quote)), escrow.totalOwed(quote));
    }

    function _seed() internal {
        LaunchHookV2.SeedPosition[] memory positions = new LaunchHookV2.SeedPosition[](1);
        positions[0] = LaunchHookV2.SeedPosition({tickLower: 60, tickUpper: 120, liquidity: 1e18});
        hook.seedLiquidity(launchKey, positions, true);
    }

    function _defaultConfig() internal returns (LaunchHookV2.PoolConfig memory) {
        return LaunchHookV2.PoolConfig({
            initialized: false,
            tokenIsCurrency0: true,
            creator: makeAddr("creator"),
            protocolTreasury: makeAddr("protocolTreasury"),
            laasTreasury: makeAddr("laasTreasury"),
            protocolFeeBps: 100,
            laasFeeBps: 50,
            creatorBps: 5000,
            platformBps: 3000,
            referrerBps: 2000,
            antiSnipeStartTotalBps: 9900,
            antiSnipeWindowSeconds: 16,
            launchTime: uint48(block.timestamp)
        });
    }
}
