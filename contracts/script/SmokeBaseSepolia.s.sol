// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console2} from "forge-std/console2.sol";

import {LaunchpadRegistry} from "../src/LaunchpadRegistry.sol";
import {MultiTenantLaunchpadFactory} from "../src/MultiTenantLaunchpadFactory.sol";

contract SmokeBaseSepolia is Script {
    function run() external returns (bytes32 launchpadId, address token, bytes32 poolId) {
        uint256 privateKey = vm.envUint("BASE_SEPOLIA_DEPLOYER_PRIVATE_KEY");
        address sender = vm.addr(privateKey);
        LaunchpadRegistry registry = LaunchpadRegistry(vm.envAddress("BASE_SEPOLIA_REGISTRY_ADDRESS"));
        MultiTenantLaunchpadFactory factory = MultiTenantLaunchpadFactory(vm.envAddress("BASE_SEPOLIA_FACTORY_ADDRESS"));
        address treasury = vm.envOr("SMOKE_LAUNCHPAD_TREASURY", sender);
        string memory slug = vm.envOr("SMOKE_LAUNCHPAD_SLUG", string("demo-pad"));
        string memory name = vm.envOr("SMOKE_TOKEN_NAME", string("O1 Demo Token"));
        string memory symbol = vm.envOr("SMOKE_TOKEN_SYMBOL", string("O1DEMO"));
        string memory contractURI = vm.envOr("SMOKE_TOKEN_URI", string("ipfs://o1-demo"));
        bytes32 salt = vm.envOr("SMOKE_LAUNCH_SALT", keccak256(abi.encode(sender, block.timestamp)));

        vm.startBroadcast(privateKey);
        launchpadId = registry.createLaunchpad(slug, treasury);
        (token, poolId) = factory.launch(
            MultiTenantLaunchpadFactory.LaunchParams({
                launchpadId: launchpadId,
                name: name,
                symbol: symbol,
                contractURI: contractURI,
                salt: salt,
                quote: factory.quote(),
                expectedConfigVersion: factory.configVersion(),
                deadline: uint64(block.timestamp + 30 minutes)
            })
        );
        vm.stopBroadcast();

        console2.logBytes32(launchpadId);
        console2.log("Token", token);
        console2.logBytes32(poolId);
    }
}
