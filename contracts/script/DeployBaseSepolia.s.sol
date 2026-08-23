// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console2} from "forge-std/console2.sol";
import {Create2} from "@openzeppelin/contracts/utils/Create2.sol";
import {SafeCast} from "@openzeppelin/contracts/utils/math/SafeCast.sol";
import {IPoolManager} from "v4-core/src/interfaces/IPoolManager.sol";
import {Hooks} from "v4-core/src/libraries/Hooks.sol";

import {FeeEscrow} from "../src/FeeEscrow.sol";
import {LaunchHookV2} from "../src/LaunchHookV2.sol";
import {LaunchpadRegistry} from "../src/LaunchpadRegistry.sol";
import {MultiTenantLaunchpadFactory} from "../src/MultiTenantLaunchpadFactory.sol";

library HookSaltMiner {
    uint160 internal constant REQUIRED_FLAGS =
        (1 << 13) | (1 << 11) | (1 << 9) | (1 << 7) | (1 << 6) | (1 << 3) | (1 << 2);
    uint160 internal constant ALL_FLAGS = (1 << 14) - 1;

    function find(address poolManager, address deployer) internal pure returns (bytes32 salt) {
        bytes32 initCodeHash = keccak256(
            abi.encodePacked(type(LaunchHookV2).creationCode, abi.encode(IPoolManager(poolManager), deployer))
        );
        for (uint256 nonce;; ++nonce) {
            salt = bytes32(nonce);
            address predicted = Create2.computeAddress(salt, initCodeHash, deployer);
            if (uint160(predicted) & ALL_FLAGS == REQUIRED_FLAGS) return salt;
        }
    }

    function predict(address poolManager, address deployer, bytes32 salt) internal pure returns (address) {
        bytes32 initCodeHash = keccak256(
            abi.encodePacked(type(LaunchHookV2).creationCode, abi.encode(IPoolManager(poolManager), deployer))
        );
        return Create2.computeAddress(salt, initCodeHash, deployer);
    }
}

contract HookCreate2Deployer {
    address public immutable owner = msg.sender;
    bool public hookDeployed;
    bool public wired;

    error AlreadyDeployed();
    error AlreadyWired();
    error HookAddressMismatch();
    error NotOwner();

    function deployHook(bytes32 salt, IPoolManager poolManager) external returns (LaunchHookV2 hook) {
        if (msg.sender != owner) revert NotOwner();
        if (hookDeployed) revert AlreadyDeployed();
        hookDeployed = true;

        address predicted = HookSaltMiner.predict(address(poolManager), address(this), salt);
        hook = new LaunchHookV2{salt: salt}(poolManager, address(this));
        if (address(hook) != predicted) revert HookAddressMismatch();
        Hooks.validateHookPermissions(hook, hook.getHookPermissions());
    }

    function wire(LaunchHookV2 hook, address factory, FeeEscrow escrow) external {
        if (msg.sender != owner) revert NotOwner();
        if (wired) revert AlreadyWired();
        wired = true;
        hook.setFactoryOnce(factory, address(escrow));
    }
}

contract DeployBaseSepolia is Script {
    using SafeCast for int256;

    struct Deployment {
        HookCreate2Deployer hookDeployer;
        LaunchpadRegistry registry;
        MultiTenantLaunchpadFactory factory;
        LaunchHookV2 hook;
        FeeEscrow escrow;
    }

    IPoolManager internal constant BASE_SEPOLIA_POOL_MANAGER = IPoolManager(0x05E73354cFDd6745C338b50BcFDfA3Aa6fA03408);

    function run() external returns (Deployment memory deployment) {
        uint256 privateKey = vm.envUint("BASE_SEPOLIA_DEPLOYER_PRIVATE_KEY");
        address quote = vm.envAddress("BASE_SEPOLIA_QUOTE_ADDRESS");
        address laasTreasury = vm.envAddress("LAAS_TREASURY_ADDRESS");
        int24 startTickToken0Frame = vm.envInt("START_TICK_TOKEN0_FRAME").toInt24();
        int24 tickSpacing = vm.envInt("TICK_SPACING").toInt24();

        vm.startBroadcast(privateKey);
        deployment.hookDeployer = new HookCreate2Deployer();
        vm.stopBroadcast();

        bytes32 hookSalt = HookSaltMiner.find(address(BASE_SEPOLIA_POOL_MANAGER), address(deployment.hookDeployer));
        address predictedHook =
            HookSaltMiner.predict(address(BASE_SEPOLIA_POOL_MANAGER), address(deployment.hookDeployer), hookSalt);

        vm.startBroadcast(privateKey);
        deployment.escrow = new FeeEscrow(BASE_SEPOLIA_POOL_MANAGER, predictedHook);
        deployment.hook = deployment.hookDeployer.deployHook(hookSalt, BASE_SEPOLIA_POOL_MANAGER);
        deployment.registry = new LaunchpadRegistry();
        deployment.factory = new MultiTenantLaunchpadFactory(
            deployment.registry,
            BASE_SEPOLIA_POOL_MANAGER,
            deployment.hook,
            quote,
            laasTreasury,
            startTickToken0Frame,
            tickSpacing
        );
        deployment.hookDeployer.wire(deployment.hook, address(deployment.factory), deployment.escrow);
        vm.stopBroadcast();

        console2.log("HookDeployer", address(deployment.hookDeployer));
        console2.log("Registry", address(deployment.registry));
        console2.log("Factory", address(deployment.factory));
        console2.log("Hook", address(deployment.hook));
        console2.log("FeeEscrow", address(deployment.escrow));
        console2.logBytes32(hookSalt);
    }
}
