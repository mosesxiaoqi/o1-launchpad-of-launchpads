// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IPoolManager} from "v4-core/src/interfaces/IPoolManager.sol";
import {IUnlockCallback} from "v4-core/src/interfaces/callback/IUnlockCallback.sol";
import {Currency} from "v4-core/src/types/Currency.sol";

import {FeeEscrow} from "../src/FeeEscrow.sol";

contract MockCurrency is ERC20 {
    constructor() ERC20("Currency", "CUR") {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract MockClaimPoolManager {
    mapping(address owner => mapping(uint256 id => uint256 amount)) public balanceOf;

    function mintClaim(address owner, address currency, uint256 amount) external {
        balanceOf[owner][uint160(currency)] += amount;
    }

    function unlock(bytes calldata data) external returns (bytes memory) {
        return IUnlockCallback(msg.sender).unlockCallback(data);
    }

    function burn(address from, uint256 id, uint256 amount) external {
        balanceOf[from][id] -= amount;
    }

    function take(Currency currency, address to, uint256 amount) external {
        address currencyAddress = Currency.unwrap(currency);
        if (currencyAddress == address(0)) {
            (bool success,) = to.call{value: amount}("");
            require(success);
        } else {
            MockCurrency(currencyAddress).transfer(to, amount);
        }
    }

    receive() external payable {}
}

contract ReentrantCurrency is MockCurrency {
    FeeEscrow internal target;
    address internal reentrantRecipient;
    address internal transferSource;
    bool public reentryAttempted;
    bool public reentrySucceeded;

    function configureReentry(FeeEscrow target_, address reentrantRecipient_, address transferSource_) external {
        target = target_;
        reentrantRecipient = reentrantRecipient_;
        transferSource = transferSource_;
    }

    function _update(address from, address to, uint256 value) internal override {
        if (address(target) != address(0) && from == transferSource && !reentryAttempted) {
            reentryAttempted = true;
            (reentrySucceeded,) =
                address(target).call(abi.encodeCall(FeeEscrow.claim, (reentrantRecipient, address(this))));
        }
        super._update(from, to, value);
    }
}

contract FeeEscrowTest is Test {
    FeeEscrow internal escrow;
    MockClaimPoolManager internal manager;
    MockCurrency internal currency;

    address internal hook = makeAddr("hook");
    address internal recipient = makeAddr("recipient");

    function setUp() public {
        manager = new MockClaimPoolManager();
        escrow = new FeeEscrow(IPoolManager(address(manager)), hook);
        currency = new MockCurrency();
    }

    function testClaimPaysAndClears() public {
        _fundClaim(currency, 10 ether);
        vm.prank(hook);
        escrow.credit(recipient, address(currency), 10 ether);

        escrow.claim(recipient, address(currency));

        assertEq(currency.balanceOf(recipient), 10 ether);
        assertEq(escrow.owed(recipient, address(currency)), 0);
    }

    function testOnlyHookCredits() public {
        _fundClaim(currency, 1 ether);

        vm.expectRevert(FeeEscrow.NotHook.selector);
        escrow.credit(recipient, address(currency), 1 ether);
    }

    function testRejectsZeroRecipientAndAmount() public {
        _fundClaim(currency, 1 ether);

        vm.startPrank(hook);
        vm.expectRevert(FeeEscrow.ZeroAddress.selector);
        escrow.credit(address(0), address(currency), 1 ether);
        vm.expectRevert(FeeEscrow.ZeroAmount.selector);
        escrow.credit(recipient, address(currency), 0);
        vm.stopPrank();
    }

    function testAggregateCreditCannotExceedActualBalance() public {
        _fundClaim(currency, 10 ether);
        vm.startPrank(hook);
        escrow.credit(recipient, address(currency), 6 ether);

        vm.expectRevert(abi.encodeWithSelector(FeeEscrow.InsufficientEscrowBalance.selector, 10 ether, 11 ether));
        escrow.credit(makeAddr("secondRecipient"), address(currency), 5 ether);
        vm.stopPrank();
    }

    function testClaimToRedirectsOnlyCallersOwnBalance() public {
        address destination = makeAddr("destination");
        _fundClaim(currency, 4 ether);
        vm.prank(hook);
        escrow.credit(recipient, address(currency), 4 ether);

        vm.prank(recipient);
        escrow.claimTo(address(currency), destination);

        assertEq(currency.balanceOf(destination), 4 ether);
        assertEq(escrow.owed(recipient, address(currency)), 0);
    }

    function testClaimPaysNativeCurrency() public {
        vm.deal(address(manager), 3 ether);
        manager.mintClaim(address(escrow), address(0), 3 ether);
        vm.prank(hook);
        escrow.credit(recipient, address(0), 3 ether);

        escrow.claim(recipient, address(0));

        assertEq(recipient.balance, 3 ether);
        assertEq(escrow.totalOwed(address(0)), 0);
    }

    function testClaimCannotPayTwice() public {
        _fundClaim(currency, 2 ether);
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
        malicious.mint(address(manager), 15 ether);
        manager.mintClaim(address(escrow), address(malicious), 15 ether);
        vm.startPrank(hook);
        escrow.credit(recipient, address(malicious), 10 ether);
        escrow.credit(secondRecipient, address(malicious), 5 ether);
        vm.stopPrank();
        malicious.configureReentry(escrow, secondRecipient, address(manager));

        escrow.claim(recipient, address(malicious));

        assertTrue(malicious.reentryAttempted());
        assertFalse(malicious.reentrySucceeded());
        assertEq(malicious.balanceOf(recipient), 10 ether);
        assertEq(escrow.owed(secondRecipient, address(malicious)), 5 ether);
    }

    function testFuzzPaidPlusOwedNeverExceedsReceived(uint128 rawAmount) public {
        uint256 amount = bound(uint256(rawAmount), 1, type(uint128).max);
        _fundClaim(currency, amount);
        vm.prank(hook);
        escrow.credit(recipient, address(currency), amount);

        assertLe(currency.balanceOf(recipient) + escrow.totalOwed(address(currency)), amount);

        escrow.claim(recipient, address(currency));

        assertEq(currency.balanceOf(recipient) + escrow.totalOwed(address(currency)), amount);
    }

    function testMultiRecipientCreditAndClaimSequencePreservesAccounting() public {
        address secondRecipient = makeAddr("secondRecipient");
        _fundClaim(currency, 30 ether);

        vm.startPrank(hook);
        escrow.credit(recipient, address(currency), 10 ether);
        _assertAccounting(secondRecipient, 10 ether, 30 ether);
        escrow.credit(secondRecipient, address(currency), 8 ether);
        _assertAccounting(secondRecipient, 18 ether, 30 ether);
        vm.stopPrank();

        escrow.claim(recipient, address(currency));
        _assertAccounting(secondRecipient, 8 ether, 20 ether);

        vm.prank(hook);
        escrow.credit(secondRecipient, address(currency), 5 ether);
        _assertAccounting(secondRecipient, 13 ether, 20 ether);

        escrow.claim(secondRecipient, address(currency));
        _assertAccounting(secondRecipient, 0, 7 ether);
        assertEq(currency.balanceOf(recipient) + currency.balanceOf(secondRecipient), 23 ether);
    }

    function _fundClaim(MockCurrency token, uint256 amount) internal {
        token.mint(address(manager), amount);
        manager.mintClaim(address(escrow), address(token), amount);
    }

    function _assertAccounting(address secondRecipient, uint256 expectedOwed, uint256 expectedBalance) internal view {
        uint256 total = escrow.owed(recipient, address(currency)) + escrow.owed(secondRecipient, address(currency));
        assertEq(total, escrow.totalOwed(address(currency)));
        assertEq(total, expectedOwed);
        assertEq(manager.balanceOf(address(escrow), uint160(address(currency))), expectedBalance);
        assertLe(total, expectedBalance);
    }
}
