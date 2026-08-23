// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

contract LaunchHookV2 {
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

    error InvalidFeeConfig();
    error PoolAlreadyRegistered();
    error UnknownPool();

    mapping(bytes32 poolId => PoolConfig config) public poolConfig;

    function currentTotalFeeBps(bytes32 poolId, uint256 timestamp) public view returns (uint256) {
        PoolConfig storage config = poolConfig[poolId];
        if (!config.initialized) revert UnknownPool();
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
        PoolConfig storage config = poolConfig[poolId];
        if (!config.initialized) revert UnknownPool();

        uint256 protocolFee = amount * config.protocolFeeBps / BPS;
        split.laas = amount * config.laasFeeBps / BPS;
        split.total = amount * currentTotalFeeBps(poolId, timestamp) / BPS;
        split.creator = protocolFee * config.creatorBps / BPS;

        bool validReferrer = referrer != address(0) && referrer != trader && referrer != config.creator
            && referrer != config.protocolTreasury && referrer != config.laasTreasury;
        split.referrer = validReferrer ? protocolFee * config.referrerBps / BPS : 0;
        split.protocol = protocolFee - split.creator - split.referrer;
        split.surcharge = split.total - protocolFee - split.laas;
    }

    function _registerPoolConfig(bytes32 poolId, PoolConfig calldata config) internal {
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
}
