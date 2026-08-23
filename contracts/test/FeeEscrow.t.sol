// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

import {FeeEscrow} from "../src/FeeEscrow.sol";

contract MockCurrency is ERC20 {
    constructor() ERC20("Currency", "CUR") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract ReentrantCurrency is MockCurrency {
    FeeEscrow internal target;
    address internal reentrantRecipient;
    bool public reentryAttempted;
    bool public reentrySucceeded;

    function configureReentry(FeeEscrow target_, address reentrantRecipient_) external {
        target = target_;
        reentrantRecipient = reentrantRecipient_;
    }

    function _update(address from, address to, uint256 value) internal override {
        if (address(target) != address(0) && from == address(target) && !reentryAttempted) {
            reentryAttempted = true;
            (reentrySucceeded,) =
                address(target).call(abi.encodeCall(FeeEscrow.claim, (reentrantRecipient, address(this))));
        }
        super._update(from, to, value);
    }
}

contract FeeEscrowTest is Test {
    FeeEscrow internal escrow;
    MockCurrency internal currency;

    address internal hook = makeAddr("hook");
    address internal recipient = makeAddr("recipient");

    function setUp() public {
        escrow = new FeeEscrow(hook);
        currency = new MockCurrency();
    }

    function testClaimPaysAndClears() public {
        currency.mint(address(escrow), 10 ether);
        vm.prank(hook);
        escrow.credit(recipient, address(currency), 10 ether);

        escrow.claim(recipient, address(currency));

        assertEq(currency.balanceOf(recipient), 10 ether);
        assertEq(escrow.owed(recipient, address(currency)), 0);
    }

    function testOnlyHookCredits() public {
        currency.mint(address(escrow), 1 ether);

        vm.expectRevert(FeeEscrow.NotHook.selector);
        escrow.credit(recipient, address(currency), 1 ether);
    }

    function testRejectsZeroRecipientAndAmount() public {
        currency.mint(address(escrow), 1 ether);

        vm.startPrank(hook);
        vm.expectRevert(FeeEscrow.ZeroAddress.selector);
        escrow.credit(address(0), address(currency), 1 ether);
        vm.expectRevert(FeeEscrow.ZeroAmount.selector);
        escrow.credit(recipient, address(currency), 0);
        vm.stopPrank();
    }

    function testAggregateCreditCannotExceedActualBalance() public {
        currency.mint(address(escrow), 10 ether);
        vm.startPrank(hook);
        escrow.credit(recipient, address(currency), 6 ether);

        vm.expectRevert(abi.encodeWithSelector(FeeEscrow.InsufficientEscrowBalance.selector, 10 ether, 11 ether));
        escrow.credit(makeAddr("secondRecipient"), address(currency), 5 ether);
        vm.stopPrank();
    }

    function testClaimToRedirectsOnlyCallersOwnBalance() public {
        address destination = makeAddr("destination");
        currency.mint(address(escrow), 4 ether);
        vm.prank(hook);
        escrow.credit(recipient, address(currency), 4 ether);

        vm.prank(recipient);
        escrow.claimTo(address(currency), destination);

        assertEq(currency.balanceOf(destination), 4 ether);
        assertEq(escrow.owed(recipient, address(currency)), 0);
    }

    function testClaimPaysNativeCurrency() public {
        vm.deal(address(escrow), 3 ether);
        vm.prank(hook);
        escrow.credit(recipient, address(0), 3 ether);

        escrow.claim(recipient, address(0));

        assertEq(recipient.balance, 3 ether);
        assertEq(escrow.totalOwed(address(0)), 0);
    }

    function testClaimCannotPayTwice() public {
        currency.mint(address(escrow), 2 ether);
        vm.prank(hook);
        escrow.credit(recipient, address(currency), 2 ether);
        escrow.claim(recipient, address(currency));

        vm.expectRevert(FeeEscrow.NothingToClaim.selector);
        escrow.claim(recipient, address(currency));
        assertEq(currency.balanceOf(recipient), 2 ether);
    }

    function testTokenTransferCannotReenterAnotherClaim() public {
        ReentrantCurrency malicious = new ReentrantCurrency();
        address secondRecipient = address(malicious);
        malicious.mint(address(escrow), 15 ether);
        vm.startPrank(hook);
        escrow.credit(recipient, address(malicious), 10 ether);
        escrow.credit(secondRecipient, address(malicious), 5 ether);
        vm.stopPrank();
        malicious.configureReentry(escrow, secondRecipient);

        escrow.claim(recipient, address(malicious));

        assertTrue(malicious.reentryAttempted());
        assertFalse(malicious.reentrySucceeded());
        assertEq(malicious.balanceOf(recipient), 10 ether);
        assertEq(escrow.owed(secondRecipient, address(malicious)), 5 ether);
    }

    function testFuzzPaidPlusOwedNeverExceedsReceived(uint128 rawAmount) public {
        uint256 amount = bound(uint256(rawAmount), 1, type(uint128).max);
        currency.mint(address(escrow), amount);
        vm.prank(hook);
        escrow.credit(recipient, address(currency), amount);

        assertLe(currency.balanceOf(recipient) + escrow.totalOwed(address(currency)), amount);

        escrow.claim(recipient, address(currency));

        assertEq(currency.balanceOf(recipient) + escrow.totalOwed(address(currency)), amount);
    }
}
