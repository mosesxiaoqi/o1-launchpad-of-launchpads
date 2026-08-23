// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Deployers} from "v4-core/test/utils/Deployers.sol";
import {Hooks} from "v4-core/src/libraries/Hooks.sol";
import {Currency} from "v4-core/src/types/Currency.sol";

import {HookCreate2Deployer, HookSaltMiner} from "../script/DeployBaseSepolia.s.sol";
import {FeeEscrow} from "../src/FeeEscrow.sol";
import {LaunchHookV2} from "../src/LaunchHookV2.sol";
import {LaunchpadRegistry} from "../src/LaunchpadRegistry.sol";
import {MultiTenantLaunchpadFactory} from "../src/MultiTenantLaunchpadFactory.sol";

contract DeploymentWiringTest is Deployers {
    address internal laasTreasury = makeAddr("laasTreasury");

    function setUp() public {
        deployFreshManagerAndRouters();
        deployMintAndApprove2Currencies();
    }

    function testDeploysSuiteWithConsistentWiringAndHookBits() public {
        HookCreate2Deployer hookDeployer = new HookCreate2Deployer();
        bytes32 hookSalt = HookSaltMiner.find(address(manager), address(hookDeployer));
        address predictedHook = HookSaltMiner.predict(address(manager), address(hookDeployer), hookSalt);
        FeeEscrow escrow = new FeeEscrow(manager, predictedHook);
        LaunchHookV2 hook = hookDeployer.deployHook(hookSalt, manager);
        LaunchpadRegistry registry = new LaunchpadRegistry();
        MultiTenantLaunchpadFactory factory =
            new MultiTenantLaunchpadFactory(registry, manager, hook, Currency.unwrap(currency1), laasTreasury, 0, 60);
        hookDeployer.wire(hook, address(factory), escrow);

        Hooks.validateHookPermissions(hook, hook.getHookPermissions());
        assertEq(address(escrow.hook()), address(hook));
        assertEq(address(escrow.poolManager()), address(manager));
        assertEq(hook.factory(), address(factory));
        assertEq(hook.feeEscrow(), address(escrow));
        assertEq(address(factory.registry()), address(registry));
        assertEq(address(factory.poolManager()), address(manager));
        assertEq(address(factory.hook()), address(hook));
    }

    function testHookDeployerCanOnlyRunOnceAndOnlyByOwner() public {
        HookCreate2Deployer hookDeployer = new HookCreate2Deployer();
        bytes32 hookSalt = HookSaltMiner.find(address(manager), address(hookDeployer));

        vm.expectRevert(HookCreate2Deployer.NotOwner.selector);
        vm.prank(makeAddr("attacker"));
        hookDeployer.deployHook(hookSalt, manager);

        hookDeployer.deployHook(hookSalt, manager);
        vm.expectRevert(HookCreate2Deployer.AlreadyDeployed.selector);
        hookDeployer.deployHook(hookSalt, manager);
    }
}
