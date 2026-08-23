// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";

import {LaunchHookV2} from "../src/LaunchHookV2.sol";

contract LaunchHookV2Harness is LaunchHookV2 {
    function registerForTest(bytes32 poolId, PoolConfig calldata config) external {
        _registerPoolConfig(poolId, config);
    }
}

contract LaunchHookV2FeesTest is Test {
    LaunchHookV2Harness internal hook;

    bytes32 internal constant POOL_ID = keccak256("pool");
    address internal creator = makeAddr("creator");
    address internal protocolTreasury = makeAddr("protocolTreasury");
    address internal laasTreasury = makeAddr("laasTreasury");
    address internal trader = makeAddr("trader");
    address internal referrer = makeAddr("referrer");
    uint48 internal constant LAUNCH_TIME = 1_700_000_000;

    function setUp() public {
        hook = new LaunchHookV2Harness();
        hook.registerForTest(POOL_ID, _defaultConfig());
    }

    function testSplitsOneThousandWithValidReferrer() public view {
        LaunchHookV2.FeeSplit memory split = hook.previewFeeSplit(POOL_ID, 1000, trader, referrer, LAUNCH_TIME + 16);

        assertEq(split.creator, 5);
        assertEq(split.protocol, 3);
        assertEq(split.referrer, 2);
        assertEq(split.laas, 5);
        assertEq(split.surcharge, 0);
        assertEq(split.total, 15);
    }

    function testInvalidReferrerRollsIntoProtocol() public view {
        LaunchHookV2.FeeSplit memory split = hook.previewFeeSplit(POOL_ID, 1000, trader, address(0), LAUNCH_TIME + 16);

        assertEq(split.creator, 5);
        assertEq(split.protocol, 5);
        assertEq(split.referrer, 0);
        assertEq(split.laas, 5);
        assertEq(split.total, 15);
    }

    function testAntiSnipeDecaysToNormalFeeAfterSixteenSeconds() public view {
        assertEq(hook.currentTotalFeeBps(POOL_ID, LAUNCH_TIME), 9900);
        assertEq(hook.currentTotalFeeBps(POOL_ID, LAUNCH_TIME + 8), 5025);
        assertEq(hook.currentTotalFeeBps(POOL_ID, LAUNCH_TIME + 16), 150);
        assertEq(hook.currentTotalFeeBps(POOL_ID, LAUNCH_TIME + 100), 150);
    }

    function testRejectsNonCanonicalProtocolSplit() public {
        LaunchHookV2.PoolConfig memory config = _defaultConfig();
        config.creatorBps = 4999;
        config.platformBps = 3001;

        vm.expectRevert(LaunchHookV2.InvalidFeeConfig.selector);
        hook.registerForTest(keccak256("other-pool"), config);
    }

    function testSelfAndTreasuryReferrersRollIntoProtocol() public view {
        address[4] memory invalid = [trader, creator, protocolTreasury, laasTreasury];
        for (uint256 i; i < invalid.length; ++i) {
            LaunchHookV2.FeeSplit memory split =
                hook.previewFeeSplit(POOL_ID, 1000, trader, invalid[i], LAUNCH_TIME + 16);
            assertEq(split.referrer, 0);
            assertEq(split.protocol, 5);
        }
    }

    function testFuzzSplitAlwaysEqualsChargedFee(uint128 amount, uint8 elapsed, bool validReferrer) public view {
        uint256 timestamp = LAUNCH_TIME + bound(uint256(elapsed), 0, 100);
        address selectedReferrer = validReferrer ? referrer : address(0);

        LaunchHookV2.FeeSplit memory split = hook.previewFeeSplit(POOL_ID, amount, trader, selectedReferrer, timestamp);

        assertEq(split.creator + split.protocol + split.referrer + split.laas + split.surcharge, split.total);
        assertEq(split.laas, uint256(amount) * 50 / 10_000);
        assertEq(split.creator, (uint256(amount) * 100 / 10_000) * 5000 / 10_000);
        assertGe(hook.currentTotalFeeBps(POOL_ID, timestamp), 150);
        assertLe(hook.currentTotalFeeBps(POOL_ID, timestamp), 9900);
    }

    function testFuzzAntiSnipeFeeIsMonotonic(uint8 firstRaw, uint8 secondRaw) public view {
        uint256 first = bound(uint256(firstRaw), 0, 16);
        uint256 second = bound(uint256(secondRaw), first, 16);

        assertGe(
            hook.currentTotalFeeBps(POOL_ID, LAUNCH_TIME + first),
            hook.currentTotalFeeBps(POOL_ID, LAUNCH_TIME + second)
        );
    }

    function _defaultConfig() internal view returns (LaunchHookV2.PoolConfig memory) {
        return LaunchHookV2.PoolConfig({
            initialized: true,
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
            launchTime: LAUNCH_TIME
        });
    }
}
