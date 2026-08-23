// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
	_ = time.Tick
	_ = context.Background
)

// FeeEscrowMetaData contains all meta data concerning the FeeEscrow contract.
var FeeEscrowMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"poolManager_\",\"type\":\"address\",\"internalType\":\"contractIPoolManager\"},{\"name\":\"hook_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claim\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"currency\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimTo\",\"inputs\":[{\"name\":\"currency\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"credit\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"currency\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hook\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owed\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"currency\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"poolManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPoolManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalOwed\",\"inputs\":[{\"name\":\"currency\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unlockCallback\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Claimed\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"currency\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Credited\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"currency\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InsufficientEscrowBalance\",\"inputs\":[{\"name\":\"available\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"required\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotHook\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotPoolManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NothingToClaim\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroDestination\",\"inputs\":[]}]",
}

// FeeEscrowABI is the input ABI used to generate the binding from.
// Deprecated: Use FeeEscrowMetaData.ABI instead.
var FeeEscrowABI = FeeEscrowMetaData.ABI

// FeeEscrow is an auto generated Go binding around an Ethereum contract.
type FeeEscrow struct {
	FeeEscrowCaller     // Read-only binding to the contract
	FeeEscrowTransactor // Write-only binding to the contract
	FeeEscrowFilterer   // Log filterer for contract events
}

// FeeEscrowCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeeEscrowCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeEscrowTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeeEscrowTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeEscrowFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeeEscrowFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeEscrowSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeeEscrowSession struct {
	Contract     *FeeEscrow        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FeeEscrowCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeeEscrowCallerSession struct {
	Contract *FeeEscrowCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// FeeEscrowTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeeEscrowTransactorSession struct {
	Contract     *FeeEscrowTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// FeeEscrowRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeeEscrowRaw struct {
	Contract *FeeEscrow // Generic contract binding to access the raw methods on
}

// FeeEscrowCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeeEscrowCallerRaw struct {
	Contract *FeeEscrowCaller // Generic read-only contract binding to access the raw methods on
}

// FeeEscrowTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeeEscrowTransactorRaw struct {
	Contract *FeeEscrowTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFeeEscrow creates a new instance of FeeEscrow, bound to a specific deployed contract.
func NewFeeEscrow(address common.Address, backend bind.ContractBackend) (*FeeEscrow, error) {
	contract, err := bindFeeEscrow(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FeeEscrow{FeeEscrowCaller: FeeEscrowCaller{contract: contract}, FeeEscrowTransactor: FeeEscrowTransactor{contract: contract}, FeeEscrowFilterer: FeeEscrowFilterer{contract: contract}}, nil
}

// NewFeeEscrowCaller creates a new read-only instance of FeeEscrow, bound to a specific deployed contract.
func NewFeeEscrowCaller(address common.Address, caller bind.ContractCaller) (*FeeEscrowCaller, error) {
	contract, err := bindFeeEscrow(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeeEscrowCaller{contract: contract}, nil
}

// NewFeeEscrowTransactor creates a new write-only instance of FeeEscrow, bound to a specific deployed contract.
func NewFeeEscrowTransactor(address common.Address, transactor bind.ContractTransactor) (*FeeEscrowTransactor, error) {
	contract, err := bindFeeEscrow(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeeEscrowTransactor{contract: contract}, nil
}

// NewFeeEscrowFilterer creates a new log filterer instance of FeeEscrow, bound to a specific deployed contract.
func NewFeeEscrowFilterer(address common.Address, filterer bind.ContractFilterer) (*FeeEscrowFilterer, error) {
	contract, err := bindFeeEscrow(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeeEscrowFilterer{contract: contract}, nil
}

// bindFeeEscrow binds a generic wrapper to an already deployed contract.
func bindFeeEscrow(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeeEscrowMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeEscrow *FeeEscrowRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeEscrow.Contract.FeeEscrowCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeEscrow *FeeEscrowRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeEscrow.Contract.FeeEscrowTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeEscrow *FeeEscrowRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeEscrow.Contract.FeeEscrowTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeEscrow *FeeEscrowCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeEscrow.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeEscrow *FeeEscrowTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeEscrow.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeEscrow *FeeEscrowTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeEscrow.Contract.contract.Transact(opts, method, params...)
}

// Hook is a free data retrieval call binding the contract method 0x7f5a7c7b.
//
// Solidity: function hook() view returns(address)
func (_FeeEscrow *FeeEscrowCaller) Hook(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeeEscrow.contract.Call(opts, &out, "hook")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Hook is a free data retrieval call binding the contract method 0x7f5a7c7b.
//
// Solidity: function hook() view returns(address)
func (_FeeEscrow *FeeEscrowSession) Hook() (common.Address, error) {
	return _FeeEscrow.Contract.Hook(&_FeeEscrow.CallOpts)
}

// Hook is a free data retrieval call binding the contract method 0x7f5a7c7b.
//
// Solidity: function hook() view returns(address)
func (_FeeEscrow *FeeEscrowCallerSession) Hook() (common.Address, error) {
	return _FeeEscrow.Contract.Hook(&_FeeEscrow.CallOpts)
}

// Owed is a free data retrieval call binding the contract method 0x28079e4a.
//
// Solidity: function owed(address recipient, address currency) view returns(uint256 amount)
func (_FeeEscrow *FeeEscrowCaller) Owed(opts *bind.CallOpts, recipient common.Address, currency common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FeeEscrow.contract.Call(opts, &out, "owed", recipient, currency)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Owed is a free data retrieval call binding the contract method 0x28079e4a.
//
// Solidity: function owed(address recipient, address currency) view returns(uint256 amount)
func (_FeeEscrow *FeeEscrowSession) Owed(recipient common.Address, currency common.Address) (*big.Int, error) {
	return _FeeEscrow.Contract.Owed(&_FeeEscrow.CallOpts, recipient, currency)
}

// Owed is a free data retrieval call binding the contract method 0x28079e4a.
//
// Solidity: function owed(address recipient, address currency) view returns(uint256 amount)
func (_FeeEscrow *FeeEscrowCallerSession) Owed(recipient common.Address, currency common.Address) (*big.Int, error) {
	return _FeeEscrow.Contract.Owed(&_FeeEscrow.CallOpts, recipient, currency)
}

// PoolManager is a free data retrieval call binding the contract method 0xdc4c90d3.
//
// Solidity: function poolManager() view returns(address)
func (_FeeEscrow *FeeEscrowCaller) PoolManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeeEscrow.contract.Call(opts, &out, "poolManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PoolManager is a free data retrieval call binding the contract method 0xdc4c90d3.
//
// Solidity: function poolManager() view returns(address)
func (_FeeEscrow *FeeEscrowSession) PoolManager() (common.Address, error) {
	return _FeeEscrow.Contract.PoolManager(&_FeeEscrow.CallOpts)
}

// PoolManager is a free data retrieval call binding the contract method 0xdc4c90d3.
//
// Solidity: function poolManager() view returns(address)
func (_FeeEscrow *FeeEscrowCallerSession) PoolManager() (common.Address, error) {
	return _FeeEscrow.Contract.PoolManager(&_FeeEscrow.CallOpts)
}

// TotalOwed is a free data retrieval call binding the contract method 0xd7e9ec04.
//
// Solidity: function totalOwed(address currency) view returns(uint256 amount)
func (_FeeEscrow *FeeEscrowCaller) TotalOwed(opts *bind.CallOpts, currency common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FeeEscrow.contract.Call(opts, &out, "totalOwed", currency)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalOwed is a free data retrieval call binding the contract method 0xd7e9ec04.
//
// Solidity: function totalOwed(address currency) view returns(uint256 amount)
func (_FeeEscrow *FeeEscrowSession) TotalOwed(currency common.Address) (*big.Int, error) {
	return _FeeEscrow.Contract.TotalOwed(&_FeeEscrow.CallOpts, currency)
}

// TotalOwed is a free data retrieval call binding the contract method 0xd7e9ec04.
//
// Solidity: function totalOwed(address currency) view returns(uint256 amount)
func (_FeeEscrow *FeeEscrowCallerSession) TotalOwed(currency common.Address) (*big.Int, error) {
	return _FeeEscrow.Contract.TotalOwed(&_FeeEscrow.CallOpts, currency)
}

// Claim is a paid mutator transaction binding the contract method 0x21c0b342.
//
// Solidity: function claim(address recipient, address currency) returns()
func (_FeeEscrow *FeeEscrowTransactor) Claim(opts *bind.TransactOpts, recipient common.Address, currency common.Address) (*types.Transaction, error) {
	return _FeeEscrow.contract.Transact(opts, "claim", recipient, currency)
}

// Claim is a paid mutator transaction binding the contract method 0x21c0b342.
//
// Solidity: function claim(address recipient, address currency) returns()
func (_FeeEscrow *FeeEscrowSession) Claim(recipient common.Address, currency common.Address) (*types.Transaction, error) {
	return _FeeEscrow.Contract.Claim(&_FeeEscrow.TransactOpts, recipient, currency)
}

// Claim is a paid mutator transaction binding the contract method 0x21c0b342.
//
// Solidity: function claim(address recipient, address currency) returns()
func (_FeeEscrow *FeeEscrowTransactorSession) Claim(recipient common.Address, currency common.Address) (*types.Transaction, error) {
	return _FeeEscrow.Contract.Claim(&_FeeEscrow.TransactOpts, recipient, currency)
}

// ClaimTo is a paid mutator transaction binding the contract method 0x34a1ca89.
//
// Solidity: function claimTo(address currency, address to) returns()
func (_FeeEscrow *FeeEscrowTransactor) ClaimTo(opts *bind.TransactOpts, currency common.Address, to common.Address) (*types.Transaction, error) {
	return _FeeEscrow.contract.Transact(opts, "claimTo", currency, to)
}

// ClaimTo is a paid mutator transaction binding the contract method 0x34a1ca89.
//
// Solidity: function claimTo(address currency, address to) returns()
func (_FeeEscrow *FeeEscrowSession) ClaimTo(currency common.Address, to common.Address) (*types.Transaction, error) {
	return _FeeEscrow.Contract.ClaimTo(&_FeeEscrow.TransactOpts, currency, to)
}

// ClaimTo is a paid mutator transaction binding the contract method 0x34a1ca89.
//
// Solidity: function claimTo(address currency, address to) returns()
func (_FeeEscrow *FeeEscrowTransactorSession) ClaimTo(currency common.Address, to common.Address) (*types.Transaction, error) {
	return _FeeEscrow.Contract.ClaimTo(&_FeeEscrow.TransactOpts, currency, to)
}

// Credit is a paid mutator transaction binding the contract method 0x7de182c5.
//
// Solidity: function credit(address recipient, address currency, uint256 amount) returns()
func (_FeeEscrow *FeeEscrowTransactor) Credit(opts *bind.TransactOpts, recipient common.Address, currency common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FeeEscrow.contract.Transact(opts, "credit", recipient, currency, amount)
}

// Credit is a paid mutator transaction binding the contract method 0x7de182c5.
//
// Solidity: function credit(address recipient, address currency, uint256 amount) returns()
func (_FeeEscrow *FeeEscrowSession) Credit(recipient common.Address, currency common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FeeEscrow.Contract.Credit(&_FeeEscrow.TransactOpts, recipient, currency, amount)
}

// Credit is a paid mutator transaction binding the contract method 0x7de182c5.
//
// Solidity: function credit(address recipient, address currency, uint256 amount) returns()
func (_FeeEscrow *FeeEscrowTransactorSession) Credit(recipient common.Address, currency common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FeeEscrow.Contract.Credit(&_FeeEscrow.TransactOpts, recipient, currency, amount)
}

// UnlockCallback is a paid mutator transaction binding the contract method 0x91dd7346.
//
// Solidity: function unlockCallback(bytes data) returns(bytes)
func (_FeeEscrow *FeeEscrowTransactor) UnlockCallback(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _FeeEscrow.contract.Transact(opts, "unlockCallback", data)
}

// UnlockCallback is a paid mutator transaction binding the contract method 0x91dd7346.
//
// Solidity: function unlockCallback(bytes data) returns(bytes)
func (_FeeEscrow *FeeEscrowSession) UnlockCallback(data []byte) (*types.Transaction, error) {
	return _FeeEscrow.Contract.UnlockCallback(&_FeeEscrow.TransactOpts, data)
}

// UnlockCallback is a paid mutator transaction binding the contract method 0x91dd7346.
//
// Solidity: function unlockCallback(bytes data) returns(bytes)
func (_FeeEscrow *FeeEscrowTransactorSession) UnlockCallback(data []byte) (*types.Transaction, error) {
	return _FeeEscrow.Contract.UnlockCallback(&_FeeEscrow.TransactOpts, data)
}

// FeeEscrowClaimedIterator is returned from FilterClaimed and is used to iterate over the raw logs and unpacked data for Claimed events raised by the FeeEscrow contract.
type FeeEscrowClaimedIterator struct {
	Event *FeeEscrowClaimed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FeeEscrowClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeEscrowClaimed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FeeEscrowClaimed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FeeEscrowClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeEscrowClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeEscrowClaimed represents a Claimed event raised by the FeeEscrow contract.
type FeeEscrowClaimed struct {
	Recipient common.Address
	Currency  common.Address
	To        common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterClaimed is a free log retrieval operation binding the contract event 0x913c992353dc81b7a8ba31496c484e9b6306bd2f6c509a649a38fdf5e1c953b2.
//
// Solidity: event Claimed(address indexed recipient, address indexed currency, address indexed to, uint256 amount)
func (_FeeEscrow *FeeEscrowFilterer) FilterClaimed(opts *bind.FilterOpts, recipient []common.Address, currency []common.Address, to []common.Address) (*FeeEscrowClaimedIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var currencyRule []interface{}
	for _, currencyItem := range currency {
		currencyRule = append(currencyRule, currencyItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeEscrow.contract.FilterLogs(opts, "Claimed", recipientRule, currencyRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FeeEscrowClaimedIterator{contract: _FeeEscrow.contract, event: "Claimed", logs: logs, sub: sub}, nil
}

// WatchClaimed is a free log subscription operation binding the contract event 0x913c992353dc81b7a8ba31496c484e9b6306bd2f6c509a649a38fdf5e1c953b2.
//
// Solidity: event Claimed(address indexed recipient, address indexed currency, address indexed to, uint256 amount)
func (_FeeEscrow *FeeEscrowFilterer) WatchClaimed(opts *bind.WatchOpts, sink chan<- *FeeEscrowClaimed, recipient []common.Address, currency []common.Address, to []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var currencyRule []interface{}
	for _, currencyItem := range currency {
		currencyRule = append(currencyRule, currencyItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeEscrow.contract.WatchLogs(opts, "Claimed", recipientRule, currencyRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeEscrowClaimed)
				if err := _FeeEscrow.contract.UnpackLog(event, "Claimed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseClaimed is a log parse operation binding the contract event 0x913c992353dc81b7a8ba31496c484e9b6306bd2f6c509a649a38fdf5e1c953b2.
//
// Solidity: event Claimed(address indexed recipient, address indexed currency, address indexed to, uint256 amount)
func (_FeeEscrow *FeeEscrowFilterer) ParseClaimed(log types.Log) (*FeeEscrowClaimed, error) {
	event := new(FeeEscrowClaimed)
	if err := _FeeEscrow.contract.UnpackLog(event, "Claimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeEscrowCreditedIterator is returned from FilterCredited and is used to iterate over the raw logs and unpacked data for Credited events raised by the FeeEscrow contract.
type FeeEscrowCreditedIterator struct {
	Event *FeeEscrowCredited // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FeeEscrowCreditedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeEscrowCredited)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FeeEscrowCredited)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FeeEscrowCreditedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeEscrowCreditedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeEscrowCredited represents a Credited event raised by the FeeEscrow contract.
type FeeEscrowCredited struct {
	Recipient common.Address
	Currency  common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCredited is a free log retrieval operation binding the contract event 0x4e45da441832cf53bdaa69235704fc0575e68210f459ee1562911024b12967d5.
//
// Solidity: event Credited(address indexed recipient, address indexed currency, uint256 amount)
func (_FeeEscrow *FeeEscrowFilterer) FilterCredited(opts *bind.FilterOpts, recipient []common.Address, currency []common.Address) (*FeeEscrowCreditedIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var currencyRule []interface{}
	for _, currencyItem := range currency {
		currencyRule = append(currencyRule, currencyItem)
	}

	logs, sub, err := _FeeEscrow.contract.FilterLogs(opts, "Credited", recipientRule, currencyRule)
	if err != nil {
		return nil, err
	}
	return &FeeEscrowCreditedIterator{contract: _FeeEscrow.contract, event: "Credited", logs: logs, sub: sub}, nil
}

// WatchCredited is a free log subscription operation binding the contract event 0x4e45da441832cf53bdaa69235704fc0575e68210f459ee1562911024b12967d5.
//
// Solidity: event Credited(address indexed recipient, address indexed currency, uint256 amount)
func (_FeeEscrow *FeeEscrowFilterer) WatchCredited(opts *bind.WatchOpts, sink chan<- *FeeEscrowCredited, recipient []common.Address, currency []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var currencyRule []interface{}
	for _, currencyItem := range currency {
		currencyRule = append(currencyRule, currencyItem)
	}

	logs, sub, err := _FeeEscrow.contract.WatchLogs(opts, "Credited", recipientRule, currencyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeEscrowCredited)
				if err := _FeeEscrow.contract.UnpackLog(event, "Credited", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCredited is a log parse operation binding the contract event 0x4e45da441832cf53bdaa69235704fc0575e68210f459ee1562911024b12967d5.
//
// Solidity: event Credited(address indexed recipient, address indexed currency, uint256 amount)
func (_FeeEscrow *FeeEscrowFilterer) ParseCredited(log types.Log) (*FeeEscrowCredited, error) {
	event := new(FeeEscrowCredited)
	if err := _FeeEscrow.contract.UnpackLog(event, "Credited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
