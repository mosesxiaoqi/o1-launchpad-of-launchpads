// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {IPoolManager} from "v4-core/src/interfaces/IPoolManager.sol";
import {IUnlockCallback} from "v4-core/src/interfaces/callback/IUnlockCallback.sol";
import {Currency} from "v4-core/src/types/Currency.sol";

contract FeeEscrow is IUnlockCallback, ReentrancyGuard {
    IPoolManager public immutable poolManager;
    address public immutable hook;

    mapping(address recipient => mapping(address currency => uint256 amount)) public owed;
    mapping(address currency => uint256 amount) public totalOwed;

    error InsufficientEscrowBalance(uint256 available, uint256 required);
    error NotHook();
    error NotPoolManager();
    error NothingToClaim();
    error ZeroAddress();
    error ZeroAmount();
    error ZeroDestination();

    event Credited(address indexed recipient, address indexed currency, uint256 amount);
    event Claimed(address indexed recipient, address indexed currency, address indexed to, uint256 amount);

    constructor(IPoolManager poolManager_, address hook_) {
        if (address(poolManager_) == address(0) || hook_ == address(0)) revert ZeroAddress();
        poolManager = poolManager_;
        hook = hook_;
    }

    function credit(address recipient, address currency, uint256 amount) external {
        if (msg.sender != hook) revert NotHook();
        if (recipient == address(0)) revert ZeroAddress();
        if (amount == 0) revert ZeroAmount();

        uint256 nextTotal = totalOwed[currency] + amount;
        uint256 balance = poolManager.balanceOf(address(this), uint160(currency));
        if (nextTotal > balance) revert InsufficientEscrowBalance(balance, nextTotal);

        owed[recipient][currency] += amount;
        totalOwed[currency] = nextTotal;
        emit Credited(recipient, currency, amount);
    }

    function claim(address recipient, address currency) external nonReentrant {
        _claim(recipient, currency, recipient);
    }

    function claimTo(address currency, address to) external nonReentrant {
        if (to == address(0)) revert ZeroDestination();
        _claim(msg.sender, currency, to);
    }

    function _claim(address recipient, address currency, address to) private {
        uint256 amount = owed[recipient][currency];
        if (amount == 0) revert NothingToClaim();
        owed[recipient][currency] = 0;
        totalOwed[currency] -= amount;
        poolManager.unlock(abi.encode(to, currency, amount));
        emit Claimed(recipient, currency, to, amount);
    }

    function unlockCallback(bytes calldata data) external returns (bytes memory) {
        if (msg.sender != address(poolManager)) revert NotPoolManager();
        (address recipient, address currencyAddress, uint256 amount) = abi.decode(data, (address, address, uint256));
        Currency currency = Currency.wrap(currencyAddress);
        poolManager.burn(address(this), currency.toId(), amount);
        poolManager.take(currency, recipient, amount);
        return "";
    }
}
