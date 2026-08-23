// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";

import {LaunchpadRegistry} from "../src/LaunchpadRegistry.sol";

contract LaunchpadRegistryTest is Test {
    event LaunchpadCreated(bytes32 indexed id, address indexed owner, address indexed treasury, string slug);

    LaunchpadRegistry internal registry;

    address internal owner = makeAddr("owner");
    address internal treasury = makeAddr("treasury");
    address internal attacker = makeAddr("attacker");

    function setUp() public {
        registry = new LaunchpadRegistry();
    }

    function testCreateStoresOwnerTreasuryAndActive() public {
        vm.warp(1_700_000_000);
        vm.prank(owner);
        bytes32 id = registry.createLaunchpad("ai-pad", treasury);

        (address storedOwner, address storedTreasury, bool active, uint64 createdAt) = registry.getLaunchpad(id);
        assertEq(storedOwner, owner);
        assertEq(storedTreasury, treasury);
        assertTrue(active);
        assertEq(createdAt, 1_700_000_000);
    }

    function testOnlyOwnerMutates() public {
        vm.prank(owner);
        bytes32 id = registry.createLaunchpad("ai-pad", treasury);

        vm.expectRevert(LaunchpadRegistry.NotOwner.selector);
        vm.prank(attacker);
        registry.setActive(id, false);
    }

    function testCreateEmitsCanonicalIdAndSlug() public {
        bytes32 expectedId = keccak256(abi.encode(block.chainid, "ai-pad"));

        vm.expectEmit(true, true, true, true);
        emit LaunchpadCreated(expectedId, owner, treasury, "ai-pad");
        vm.prank(owner);
        registry.createLaunchpad("ai-pad", treasury);
    }

    function testRejectsDuplicateId() public {
        vm.prank(owner);
        registry.createLaunchpad("ai-pad", treasury);

        vm.expectRevert(LaunchpadRegistry.AlreadyExists.selector);
        registry.createLaunchpad("ai-pad", treasury);
    }

    function testRejectsInvalidSlugs() public {
        string[7] memory invalid =
            [string("ab"), "abcdefghijklmnopqrstuvwxyz1234567", "-abc", "abc-", "a--b", "Aaa", "a_b"];

        for (uint256 i; i < invalid.length; ++i) {
            vm.expectRevert(LaunchpadRegistry.InvalidSlug.selector);
            registry.createLaunchpad(invalid[i], treasury);
        }
    }

    function testRejectsZeroTreasury() public {
        vm.expectRevert(LaunchpadRegistry.ZeroTreasury.selector);
        registry.createLaunchpad("ai-pad", address(0));
    }

    function testOwnerUpdatesTreasury() public {
        vm.prank(owner);
        bytes32 id = registry.createLaunchpad("ai-pad", treasury);
        address nextTreasury = makeAddr("nextTreasury");

        vm.prank(owner);
        registry.setTreasury(id, nextTreasury);

        (, address storedTreasury,,) = registry.getLaunchpad(id);
        assertEq(storedTreasury, nextTreasury);
    }

    function testOwnerUpdatesActiveState() public {
        vm.prank(owner);
        bytes32 id = registry.createLaunchpad("ai-pad", treasury);

        vm.prank(owner);
        registry.setActive(id, false);

        (,, bool active,) = registry.getLaunchpad(id);
        assertFalse(active);
    }

    function testRejectsZeroTreasuryUpdate() public {
        vm.prank(owner);
        bytes32 id = registry.createLaunchpad("ai-pad", treasury);

        vm.expectRevert(LaunchpadRegistry.ZeroTreasury.selector);
        vm.prank(owner);
        registry.setTreasury(id, address(0));
    }
}
