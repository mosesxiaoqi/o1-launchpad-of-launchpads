// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {IHooks} from "v4-core/src/interfaces/IHooks.sol";
import {IPoolManager} from "v4-core/src/interfaces/IPoolManager.sol";
import {PoolKey} from "v4-core/src/types/PoolKey.sol";
import {Currency} from "v4-core/src/types/Currency.sol";
import {toBalanceDelta} from "v4-core/src/types/BalanceDelta.sol";
import {ModifyLiquidityParams} from "v4-core/src/types/PoolOperation.sol";
import {SwapParams} from "v4-core/src/types/PoolOperation.sol";

import {FeeEscrow} from "../src/FeeEscrow.sol";
import {LaunchHookV2} from "../src/LaunchHookV2.sol";
import {BaseHook} from "../src/base/BaseHook.sol";

contract HookPoolManagerCaller {
    function beforeInitialize(IHooks hook, address sender, PoolKey calldata key) external returns (bytes4) {
        return hook.beforeInitialize(sender, key, 1 << 96);
    }

    function beforeAddLiquidity(IHooks hook, PoolKey calldata key) external returns (bytes4) {
        return hook.beforeAddLiquidity(
            msg.sender,
            key,
            ModifyLiquidityParams({tickLower: -120, tickUpper: 120, liquidityDelta: 1, salt: bytes32(0)}),
            ""
        );
    }

    function beforeRemoveLiquidity(IHooks hook, PoolKey calldata key) external returns (bytes4) {
        return hook.beforeRemoveLiquidity(
            msg.sender,
            key,
            ModifyLiquidityParams({tickLower: -120, tickUpper: 120, liquidityDelta: -1, salt: bytes32(0)}),
            ""
        );
    }

    function beforeExactOutputSwap(IHooks hook, PoolKey calldata key) external {
        hook.beforeSwap(
            msg.sender,
            key,
            SwapParams({zeroForOne: false, amountSpecified: 1, sqrtPriceLimitX96: type(uint160).max - 1}),
            ""
        );
    }

    function afterPartialExactInputSwap(IHooks hook, PoolKey calldata key) external {
        hook.afterSwap(
            msg.sender,
            key,
            SwapParams({zeroForOne: false, amountSpecified: -1000, sqrtPriceLimitX96: type(uint160).max - 1}),
            toBalanceDelta(0, -984),
            ""
        );
    }
}

contract LaunchHookV2PoolTest is Test {
    HookPoolManagerCaller internal manager;
    LaunchHookV2 internal hook;
    FeeEscrow internal escrow;
    PoolKey internal key;

    address internal factory = makeAddr("factory");
    address internal creator = makeAddr("creator");
    address internal protocolTreasury = makeAddr("protocolTreasury");
    address internal laasTreasury = makeAddr("laasTreasury");

    function setUp() public {
        manager = new HookPoolManagerCaller();
        hook = new LaunchHookV2(IPoolManager(address(manager)), address(this));
        escrow = new FeeEscrow(IPoolManager(address(manager)), address(hook));
        key = PoolKey({
            currency0: Currency.wrap(address(1)),
            currency1: Currency.wrap(address(2)),
            fee: 3000,
            tickSpacing: 60,
            hooks: IHooks(address(hook))
        });
    }

    function testWiringCanOnlyHappenOnce() public {
        hook.setFactoryOnce(factory, address(escrow));

        vm.expectRevert(LaunchHookV2.AlreadyWired.selector);
        hook.setFactoryOnce(factory, address(escrow));
    }

    function testWiringRejectsEscrowUsingAnotherPoolManager() public {
        HookPoolManagerCaller otherManager = new HookPoolManagerCaller();
        FeeEscrow wrongEscrow = new FeeEscrow(IPoolManager(address(otherManager)), address(hook));

        vm.expectRevert(LaunchHookV2.BadWiring.selector);
        hook.setFactoryOnce(factory, address(wrongEscrow));
    }

    function testOnlyPoolManagerCallsCallbacks() public {
        vm.expectRevert(BaseHook.NotPoolManager.selector);
        hook.beforeInitialize(factory, key, 1 << 96);
    }

    function testPoolCannotInitializeBeforeWiring() public {
        vm.expectRevert(LaunchHookV2.NotWired.selector);
        manager.beforeInitialize(IHooks(address(hook)), factory, key);
    }

    function testOnlyFactoryRegistersPoolOnce() public {
        hook.setFactoryOnce(factory, address(escrow));
        LaunchHookV2.PoolConfig memory config = _defaultConfig();

        vm.expectRevert(LaunchHookV2.NotFactory.selector);
        hook.registerPool(key, config);

        vm.prank(factory);
        hook.registerPool(key, config);
        vm.expectRevert(LaunchHookV2.PoolAlreadyRegistered.selector);
        vm.prank(factory);
        hook.registerPool(key, config);
    }

    function testExternalLiquidityChangesAreLocked() public {
        hook.setFactoryOnce(factory, address(escrow));

        vm.expectRevert(LaunchHookV2.LiquidityLocked.selector);
        manager.beforeAddLiquidity(IHooks(address(hook)), key);
        vm.expectRevert(LaunchHookV2.LiquidityLocked.selector);
        manager.beforeRemoveLiquidity(IHooks(address(hook)), key);
    }

    function testExactOutputIsDisabledDuringAntiSnipe() public {
        hook.setFactoryOnce(factory, address(escrow));
        vm.prank(factory);
        hook.registerPool(key, _defaultConfig());

        vm.expectRevert(LaunchHookV2.ExactOutputDisabledDuringAntiSnipe.selector);
        manager.beforeExactOutputSwap(IHooks(address(hook)), key);

        vm.warp(block.timestamp + 16);
        manager.beforeExactOutputSwap(IHooks(address(hook)), key);
    }

    function testRejectsPartialSpecifiedFill() public {
        hook.setFactoryOnce(factory, address(escrow));
        vm.prank(factory);
        hook.registerPool(key, _defaultConfig());
        vm.warp(block.timestamp + 16);

        vm.expectRevert(LaunchHookV2.PartialFillUnsupported.selector);
        manager.afterPartialExactInputSwap(IHooks(address(hook)), key);
    }

    function _defaultConfig() internal view returns (LaunchHookV2.PoolConfig memory) {
        return LaunchHookV2.PoolConfig({
            initialized: false,
            tokenIsCurrency0: true,
            creator: creator,
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
        });
    }
}
