// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

contract LaunchpadRegistry {
    struct Launchpad {
        address owner;
        address treasury;
        bool active;
        uint64 createdAt;
    }

    error AlreadyExists();
    error InvalidSlug();
    error NotOwner();
    error ZeroTreasury();

    event LaunchpadCreated(bytes32 indexed id, address indexed owner, address indexed treasury, string slug);

    mapping(bytes32 id => Launchpad launchpad) private launchpads;

    function deriveId(string memory normalizedSlug) public view returns (bytes32) {
        return keccak256(abi.encode(block.chainid, normalizedSlug));
    }

    function createLaunchpad(string calldata normalizedSlug, address treasury) external returns (bytes32 id) {
        if (!_validSlug(bytes(normalizedSlug))) revert InvalidSlug();
        if (treasury == address(0)) revert ZeroTreasury();

        id = deriveId(normalizedSlug);
        if (launchpads[id].owner != address(0)) revert AlreadyExists();

        launchpads[id] = Launchpad(msg.sender, treasury, true, uint64(block.timestamp));
        emit LaunchpadCreated(id, msg.sender, treasury, normalizedSlug);
    }

    function getLaunchpad(bytes32 id)
        external
        view
        returns (address owner, address treasury, bool active, uint64 createdAt)
    {
        Launchpad storage launchpad = launchpads[id];
        return (launchpad.owner, launchpad.treasury, launchpad.active, launchpad.createdAt);
    }

    function setTreasury(bytes32 id, address treasury) external {
        Launchpad storage launchpad = launchpads[id];
        if (launchpad.owner != msg.sender) revert NotOwner();
        if (treasury == address(0)) revert ZeroTreasury();
        launchpad.treasury = treasury;
    }

    function setActive(bytes32 id, bool active) external {
        Launchpad storage launchpad = launchpads[id];
        if (launchpad.owner != msg.sender) revert NotOwner();
        launchpad.active = active;
    }

    function _validSlug(bytes memory slug) private pure returns (bool) {
        uint256 length = slug.length;
        if (length < 3 || length > 32 || slug[0] == "-" || slug[length - 1] == "-") return false;

        bool previousWasHyphen;
        for (uint256 i; i < length; ++i) {
            bytes1 character = slug[i];
            bool isHyphen = character == "-";
            bool isLowercase = character >= "a" && character <= "z";
            bool isDigit = character >= "0" && character <= "9";
            if ((!isHyphen && !isLowercase && !isDigit) || (isHyphen && previousWasHyphen)) return false;
            previousWasHyphen = isHyphen;
        }
        return true;
    }
}
