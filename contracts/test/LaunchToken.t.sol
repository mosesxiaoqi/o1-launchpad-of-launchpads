// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";

import {LaunchToken} from "../src/LaunchToken.sol";

contract LaunchTokenTest is Test {
    address internal recipient = makeAddr("recipient");

    function testMintsCompleteSupplyOnce() public {
        LaunchToken token = new LaunchToken("Alpha", "ALPHA", "ipfs://meta", 1_000_000_000 ether, recipient);

        assertEq(token.name(), "Alpha");
        assertEq(token.symbol(), "ALPHA");
        assertEq(token.contractURI(), "ipfs://meta");
        assertEq(token.totalSupply(), 1_000_000_000 ether);
        assertEq(token.balanceOf(recipient), token.totalSupply());
    }

    function testHasNoMintAuthorityAfterConstruction() public {
        LaunchToken token = new LaunchToken("Alpha", "ALPHA", "ipfs://meta", 1 ether, recipient);

        (bool success,) = address(token).call(abi.encodeWithSignature("mint(address,uint256)", recipient, 1 ether));

        assertFalse(success);
        assertEq(token.totalSupply(), 1 ether);
    }
}
