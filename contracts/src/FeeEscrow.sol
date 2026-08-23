// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

contract FeeEscrow is ReentrancyGuard {
    using SafeERC20 for IERC20;

    address public immutable hook;

    mapping(address recipient => mapping(address currency => uint256 amount)) public owed;
    mapping(address currency => uint256 amount) public totalOwed;

    error InsufficientEscrowBalance(uint256 available, uint256 required);
    error NativeTransferFailed();
    error NotHook();
    error NothingToClaim();
    error ZeroAddress();
    error ZeroAmount();
    error ZeroDestination();

    event Credited(address indexed recipient, address indexed currency, uint256 amount);
    event Claimed(address indexed recipient, address indexed currency, address indexed to, uint256 amount);

    constructor(address hook_) {
        if (hook_ == address(0)) revert ZeroAddress();
        hook = hook_;
    }

    receive() external payable {}

    function credit(address recipient, address currency, uint256 amount) external {
        if (msg.sender != hook) revert NotHook();
        if (recipient == address(0)) revert ZeroAddress();
        if (amount == 0) revert ZeroAmount();

        uint256 nextTotal = totalOwed[currency] + amount;
        uint256 balance = currency == address(0) ? address(this).balance : IERC20(currency).balanceOf(address(this));
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

        if (currency == address(0)) {
            (bool success,) = to.call{value: amount}("");
            if (!success) revert NativeTransferFailed();
        } else {
            IERC20(currency).safeTransfer(to, amount);
        }

        emit Claimed(recipient, currency, to, amount);
    }
}
