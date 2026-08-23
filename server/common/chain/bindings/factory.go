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

// MultiTenantLaunchpadFactoryLaunchParams is an auto generated low-level Go binding around an user-defined struct.
type MultiTenantLaunchpadFactoryLaunchParams struct {
	LaunchpadId           [32]byte
	Name                  string
	Symbol                string
	ContractURI           string
	Salt                  [32]byte
	Quote                 common.Address
	ExpectedConfigVersion uint64
	Deadline              uint64
}

// MultiTenantLaunchpadFactoryMetaData contains all meta data concerning the MultiTenantLaunchpadFactory contract.
var MultiTenantLaunchpadFactoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"registry_\",\"type\":\"address\",\"internalType\":\"contractLaunchpadRegistry\"},{\"name\":\"poolManager_\",\"type\":\"address\",\"internalType\":\"contractIPoolManager\"},{\"name\":\"hook_\",\"type\":\"address\",\"internalType\":\"contractLaunchHookV2\"},{\"name\":\"quote_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"laasTreasury_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"startTickToken0Frame_\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"tickSpacing_\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"LAUNCH_SUPPLY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"configVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hook\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractLaunchHookV2\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"laasTreasury\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"launch\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structMultiTenantLaunchpadFactory.LaunchParams\",\"components\":[{\"name\":\"launchpadId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"contractURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quote\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"expectedConfigVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"deadline\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"outputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"poolId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"launchpadOf\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"launchpadId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"poolManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPoolManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"poolOf\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"poolId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quote\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractLaunchpadRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"startTickToken0Frame\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tickSpacing\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"usedLaunchSalts\",\"inputs\":[{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"used\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Launched\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"poolId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"launchpadId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"creator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"quote\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"supply\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"tickSpacing\",\"type\":\"int24\",\"indexed\":false,\"internalType\":\"int24\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"EmptyName\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidLaunchpad\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LaunchExpired\",\"inputs\":[{\"name\":\"deadline\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"LaunchSaltUsed\",\"inputs\":[{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StaleConfig\",\"inputs\":[{\"name\":\"expectedVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"actualVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}]",
}

// MultiTenantLaunchpadFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use MultiTenantLaunchpadFactoryMetaData.ABI instead.
var MultiTenantLaunchpadFactoryABI = MultiTenantLaunchpadFactoryMetaData.ABI

// MultiTenantLaunchpadFactory is an auto generated Go binding around an Ethereum contract.
type MultiTenantLaunchpadFactory struct {
	MultiTenantLaunchpadFactoryCaller     // Read-only binding to the contract
	MultiTenantLaunchpadFactoryTransactor // Write-only binding to the contract
	MultiTenantLaunchpadFactoryFilterer   // Log filterer for contract events
}

// MultiTenantLaunchpadFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type MultiTenantLaunchpadFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiTenantLaunchpadFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MultiTenantLaunchpadFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiTenantLaunchpadFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MultiTenantLaunchpadFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiTenantLaunchpadFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MultiTenantLaunchpadFactorySession struct {
	Contract     *MultiTenantLaunchpadFactory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                // Call options to use throughout this session
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// MultiTenantLaunchpadFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MultiTenantLaunchpadFactoryCallerSession struct {
	Contract *MultiTenantLaunchpadFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                      // Call options to use throughout this session
}

// MultiTenantLaunchpadFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MultiTenantLaunchpadFactoryTransactorSession struct {
	Contract     *MultiTenantLaunchpadFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                      // Transaction auth options to use throughout this session
}

// MultiTenantLaunchpadFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type MultiTenantLaunchpadFactoryRaw struct {
	Contract *MultiTenantLaunchpadFactory // Generic contract binding to access the raw methods on
}

// MultiTenantLaunchpadFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MultiTenantLaunchpadFactoryCallerRaw struct {
	Contract *MultiTenantLaunchpadFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// MultiTenantLaunchpadFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MultiTenantLaunchpadFactoryTransactorRaw struct {
	Contract *MultiTenantLaunchpadFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMultiTenantLaunchpadFactory creates a new instance of MultiTenantLaunchpadFactory, bound to a specific deployed contract.
func NewMultiTenantLaunchpadFactory(address common.Address, backend bind.ContractBackend) (*MultiTenantLaunchpadFactory, error) {
	contract, err := bindMultiTenantLaunchpadFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MultiTenantLaunchpadFactory{MultiTenantLaunchpadFactoryCaller: MultiTenantLaunchpadFactoryCaller{contract: contract}, MultiTenantLaunchpadFactoryTransactor: MultiTenantLaunchpadFactoryTransactor{contract: contract}, MultiTenantLaunchpadFactoryFilterer: MultiTenantLaunchpadFactoryFilterer{contract: contract}}, nil
}

// NewMultiTenantLaunchpadFactoryCaller creates a new read-only instance of MultiTenantLaunchpadFactory, bound to a specific deployed contract.
func NewMultiTenantLaunchpadFactoryCaller(address common.Address, caller bind.ContractCaller) (*MultiTenantLaunchpadFactoryCaller, error) {
	contract, err := bindMultiTenantLaunchpadFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultiTenantLaunchpadFactoryCaller{contract: contract}, nil
}

// NewMultiTenantLaunchpadFactoryTransactor creates a new write-only instance of MultiTenantLaunchpadFactory, bound to a specific deployed contract.
func NewMultiTenantLaunchpadFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*MultiTenantLaunchpadFactoryTransactor, error) {
	contract, err := bindMultiTenantLaunchpadFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultiTenantLaunchpadFactoryTransactor{contract: contract}, nil
}

// NewMultiTenantLaunchpadFactoryFilterer creates a new log filterer instance of MultiTenantLaunchpadFactory, bound to a specific deployed contract.
func NewMultiTenantLaunchpadFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*MultiTenantLaunchpadFactoryFilterer, error) {
	contract, err := bindMultiTenantLaunchpadFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultiTenantLaunchpadFactoryFilterer{contract: contract}, nil
}

// bindMultiTenantLaunchpadFactory binds a generic wrapper to an already deployed contract.
func bindMultiTenantLaunchpadFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MultiTenantLaunchpadFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiTenantLaunchpadFactory.Contract.MultiTenantLaunchpadFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiTenantLaunchpadFactory.Contract.MultiTenantLaunchpadFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiTenantLaunchpadFactory.Contract.MultiTenantLaunchpadFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiTenantLaunchpadFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiTenantLaunchpadFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiTenantLaunchpadFactory.Contract.contract.Transact(opts, method, params...)
}

// LAUNCHSUPPLY is a free data retrieval call binding the contract method 0xfb1ace18.
//
// Solidity: function LAUNCH_SUPPLY() view returns(uint256)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCaller) LAUNCHSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MultiTenantLaunchpadFactory.contract.Call(opts, &out, "LAUNCH_SUPPLY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LAUNCHSUPPLY is a free data retrieval call binding the contract method 0xfb1ace18.
//
// Solidity: function LAUNCH_SUPPLY() view returns(uint256)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactorySession) LAUNCHSUPPLY() (*big.Int, error) {
	return _MultiTenantLaunchpadFactory.Contract.LAUNCHSUPPLY(&_MultiTenantLaunchpadFactory.CallOpts)
}

// LAUNCHSUPPLY is a free data retrieval call binding the contract method 0xfb1ace18.
//
// Solidity: function LAUNCH_SUPPLY() view returns(uint256)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCallerSession) LAUNCHSUPPLY() (*big.Int, error) {
	return _MultiTenantLaunchpadFactory.Contract.LAUNCHSUPPLY(&_MultiTenantLaunchpadFactory.CallOpts)
}

// ConfigVersion is a free data retrieval call binding the contract method 0xdd64d24d.
//
// Solidity: function configVersion() view returns(uint64)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCaller) ConfigVersion(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _MultiTenantLaunchpadFactory.contract.Call(opts, &out, "configVersion")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ConfigVersion is a free data retrieval call binding the contract method 0xdd64d24d.
//
// Solidity: function configVersion() view returns(uint64)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactorySession) ConfigVersion() (uint64, error) {
	return _MultiTenantLaunchpadFactory.Contract.ConfigVersion(&_MultiTenantLaunchpadFactory.CallOpts)
}

// ConfigVersion is a free data retrieval call binding the contract method 0xdd64d24d.
//
// Solidity: function configVersion() view returns(uint64)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCallerSession) ConfigVersion() (uint64, error) {
	return _MultiTenantLaunchpadFactory.Contract.ConfigVersion(&_MultiTenantLaunchpadFactory.CallOpts)
}

// Hook is a free data retrieval call binding the contract method 0x7f5a7c7b.
//
// Solidity: function hook() view returns(address)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCaller) Hook(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MultiTenantLaunchpadFactory.contract.Call(opts, &out, "hook")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Hook is a free data retrieval call binding the contract method 0x7f5a7c7b.
//
// Solidity: function hook() view returns(address)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactorySession) Hook() (common.Address, error) {
	return _MultiTenantLaunchpadFactory.Contract.Hook(&_MultiTenantLaunchpadFactory.CallOpts)
}

// Hook is a free data retrieval call binding the contract method 0x7f5a7c7b.
//
// Solidity: function hook() view returns(address)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCallerSession) Hook() (common.Address, error) {
	return _MultiTenantLaunchpadFactory.Contract.Hook(&_MultiTenantLaunchpadFactory.CallOpts)
}

// LaasTreasury is a free data retrieval call binding the contract method 0x078c15a9.
//
// Solidity: function laasTreasury() view returns(address)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCaller) LaasTreasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MultiTenantLaunchpadFactory.contract.Call(opts, &out, "laasTreasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LaasTreasury is a free data retrieval call binding the contract method 0x078c15a9.
//
// Solidity: function laasTreasury() view returns(address)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactorySession) LaasTreasury() (common.Address, error) {
	return _MultiTenantLaunchpadFactory.Contract.LaasTreasury(&_MultiTenantLaunchpadFactory.CallOpts)
}

// LaasTreasury is a free data retrieval call binding the contract method 0x078c15a9.
//
// Solidity: function laasTreasury() view returns(address)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCallerSession) LaasTreasury() (common.Address, error) {
	return _MultiTenantLaunchpadFactory.Contract.LaasTreasury(&_MultiTenantLaunchpadFactory.CallOpts)
}

// LaunchpadOf is a free data retrieval call binding the contract method 0x642d0ce2.
//
// Solidity: function launchpadOf(address token) view returns(bytes32 launchpadId)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCaller) LaunchpadOf(opts *bind.CallOpts, token common.Address) ([32]byte, error) {
	var out []interface{}
	err := _MultiTenantLaunchpadFactory.contract.Call(opts, &out, "launchpadOf", token)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LaunchpadOf is a free data retrieval call binding the contract method 0x642d0ce2.
//
// Solidity: function launchpadOf(address token) view returns(bytes32 launchpadId)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactorySession) LaunchpadOf(token common.Address) ([32]byte, error) {
	return _MultiTenantLaunchpadFactory.Contract.LaunchpadOf(&_MultiTenantLaunchpadFactory.CallOpts, token)
}

// LaunchpadOf is a free data retrieval call binding the contract method 0x642d0ce2.
//
// Solidity: function launchpadOf(address token) view returns(bytes32 launchpadId)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCallerSession) LaunchpadOf(token common.Address) ([32]byte, error) {
	return _MultiTenantLaunchpadFactory.Contract.LaunchpadOf(&_MultiTenantLaunchpadFactory.CallOpts, token)
}

// PoolManager is a free data retrieval call binding the contract method 0xdc4c90d3.
//
// Solidity: function poolManager() view returns(address)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCaller) PoolManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MultiTenantLaunchpadFactory.contract.Call(opts, &out, "poolManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PoolManager is a free data retrieval call binding the contract method 0xdc4c90d3.
//
// Solidity: function poolManager() view returns(address)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactorySession) PoolManager() (common.Address, error) {
	return _MultiTenantLaunchpadFactory.Contract.PoolManager(&_MultiTenantLaunchpadFactory.CallOpts)
}

// PoolManager is a free data retrieval call binding the contract method 0xdc4c90d3.
//
// Solidity: function poolManager() view returns(address)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCallerSession) PoolManager() (common.Address, error) {
	return _MultiTenantLaunchpadFactory.Contract.PoolManager(&_MultiTenantLaunchpadFactory.CallOpts)
}

// PoolOf is a free data retrieval call binding the contract method 0x988b1fa7.
//
// Solidity: function poolOf(address token) view returns(bytes32 poolId)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCaller) PoolOf(opts *bind.CallOpts, token common.Address) ([32]byte, error) {
	var out []interface{}
	err := _MultiTenantLaunchpadFactory.contract.Call(opts, &out, "poolOf", token)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PoolOf is a free data retrieval call binding the contract method 0x988b1fa7.
//
// Solidity: function poolOf(address token) view returns(bytes32 poolId)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactorySession) PoolOf(token common.Address) ([32]byte, error) {
	return _MultiTenantLaunchpadFactory.Contract.PoolOf(&_MultiTenantLaunchpadFactory.CallOpts, token)
}

// PoolOf is a free data retrieval call binding the contract method 0x988b1fa7.
//
// Solidity: function poolOf(address token) view returns(bytes32 poolId)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCallerSession) PoolOf(token common.Address) ([32]byte, error) {
	return _MultiTenantLaunchpadFactory.Contract.PoolOf(&_MultiTenantLaunchpadFactory.CallOpts, token)
}

// Quote is a free data retrieval call binding the contract method 0x999b93af.
//
// Solidity: function quote() view returns(address)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCaller) Quote(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MultiTenantLaunchpadFactory.contract.Call(opts, &out, "quote")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Quote is a free data retrieval call binding the contract method 0x999b93af.
//
// Solidity: function quote() view returns(address)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactorySession) Quote() (common.Address, error) {
	return _MultiTenantLaunchpadFactory.Contract.Quote(&_MultiTenantLaunchpadFactory.CallOpts)
}

// Quote is a free data retrieval call binding the contract method 0x999b93af.
//
// Solidity: function quote() view returns(address)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCallerSession) Quote() (common.Address, error) {
	return _MultiTenantLaunchpadFactory.Contract.Quote(&_MultiTenantLaunchpadFactory.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCaller) Registry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MultiTenantLaunchpadFactory.contract.Call(opts, &out, "registry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactorySession) Registry() (common.Address, error) {
	return _MultiTenantLaunchpadFactory.Contract.Registry(&_MultiTenantLaunchpadFactory.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCallerSession) Registry() (common.Address, error) {
	return _MultiTenantLaunchpadFactory.Contract.Registry(&_MultiTenantLaunchpadFactory.CallOpts)
}

// StartTickToken0Frame is a free data retrieval call binding the contract method 0x9da0360b.
//
// Solidity: function startTickToken0Frame() view returns(int24)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCaller) StartTickToken0Frame(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MultiTenantLaunchpadFactory.contract.Call(opts, &out, "startTickToken0Frame")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StartTickToken0Frame is a free data retrieval call binding the contract method 0x9da0360b.
//
// Solidity: function startTickToken0Frame() view returns(int24)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactorySession) StartTickToken0Frame() (*big.Int, error) {
	return _MultiTenantLaunchpadFactory.Contract.StartTickToken0Frame(&_MultiTenantLaunchpadFactory.CallOpts)
}

// StartTickToken0Frame is a free data retrieval call binding the contract method 0x9da0360b.
//
// Solidity: function startTickToken0Frame() view returns(int24)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCallerSession) StartTickToken0Frame() (*big.Int, error) {
	return _MultiTenantLaunchpadFactory.Contract.StartTickToken0Frame(&_MultiTenantLaunchpadFactory.CallOpts)
}

// TickSpacing is a free data retrieval call binding the contract method 0xd0c93a7c.
//
// Solidity: function tickSpacing() view returns(int24)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCaller) TickSpacing(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MultiTenantLaunchpadFactory.contract.Call(opts, &out, "tickSpacing")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TickSpacing is a free data retrieval call binding the contract method 0xd0c93a7c.
//
// Solidity: function tickSpacing() view returns(int24)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactorySession) TickSpacing() (*big.Int, error) {
	return _MultiTenantLaunchpadFactory.Contract.TickSpacing(&_MultiTenantLaunchpadFactory.CallOpts)
}

// TickSpacing is a free data retrieval call binding the contract method 0xd0c93a7c.
//
// Solidity: function tickSpacing() view returns(int24)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCallerSession) TickSpacing() (*big.Int, error) {
	return _MultiTenantLaunchpadFactory.Contract.TickSpacing(&_MultiTenantLaunchpadFactory.CallOpts)
}

// UsedLaunchSalts is a free data retrieval call binding the contract method 0x125c9280.
//
// Solidity: function usedLaunchSalts(bytes32 salt) view returns(bool used)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCaller) UsedLaunchSalts(opts *bind.CallOpts, salt [32]byte) (bool, error) {
	var out []interface{}
	err := _MultiTenantLaunchpadFactory.contract.Call(opts, &out, "usedLaunchSalts", salt)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UsedLaunchSalts is a free data retrieval call binding the contract method 0x125c9280.
//
// Solidity: function usedLaunchSalts(bytes32 salt) view returns(bool used)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactorySession) UsedLaunchSalts(salt [32]byte) (bool, error) {
	return _MultiTenantLaunchpadFactory.Contract.UsedLaunchSalts(&_MultiTenantLaunchpadFactory.CallOpts, salt)
}

// UsedLaunchSalts is a free data retrieval call binding the contract method 0x125c9280.
//
// Solidity: function usedLaunchSalts(bytes32 salt) view returns(bool used)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryCallerSession) UsedLaunchSalts(salt [32]byte) (bool, error) {
	return _MultiTenantLaunchpadFactory.Contract.UsedLaunchSalts(&_MultiTenantLaunchpadFactory.CallOpts, salt)
}

// Launch is a paid mutator transaction binding the contract method 0x22d8980f.
//
// Solidity: function launch((bytes32,string,string,string,bytes32,address,uint64,uint64) params) returns(address token, bytes32 poolId)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryTransactor) Launch(opts *bind.TransactOpts, params MultiTenantLaunchpadFactoryLaunchParams) (*types.Transaction, error) {
	return _MultiTenantLaunchpadFactory.contract.Transact(opts, "launch", params)
}

// Launch is a paid mutator transaction binding the contract method 0x22d8980f.
//
// Solidity: function launch((bytes32,string,string,string,bytes32,address,uint64,uint64) params) returns(address token, bytes32 poolId)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactorySession) Launch(params MultiTenantLaunchpadFactoryLaunchParams) (*types.Transaction, error) {
	return _MultiTenantLaunchpadFactory.Contract.Launch(&_MultiTenantLaunchpadFactory.TransactOpts, params)
}

// Launch is a paid mutator transaction binding the contract method 0x22d8980f.
//
// Solidity: function launch((bytes32,string,string,string,bytes32,address,uint64,uint64) params) returns(address token, bytes32 poolId)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryTransactorSession) Launch(params MultiTenantLaunchpadFactoryLaunchParams) (*types.Transaction, error) {
	return _MultiTenantLaunchpadFactory.Contract.Launch(&_MultiTenantLaunchpadFactory.TransactOpts, params)
}

// MultiTenantLaunchpadFactoryLaunchedIterator is returned from FilterLaunched and is used to iterate over the raw logs and unpacked data for Launched events raised by the MultiTenantLaunchpadFactory contract.
type MultiTenantLaunchpadFactoryLaunchedIterator struct {
	Event *MultiTenantLaunchpadFactoryLaunched // Event containing the contract specifics and raw log

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
func (it *MultiTenantLaunchpadFactoryLaunchedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultiTenantLaunchpadFactoryLaunched)
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
		it.Event = new(MultiTenantLaunchpadFactoryLaunched)
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
func (it *MultiTenantLaunchpadFactoryLaunchedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultiTenantLaunchpadFactoryLaunchedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultiTenantLaunchpadFactoryLaunched represents a Launched event raised by the MultiTenantLaunchpadFactory contract.
type MultiTenantLaunchpadFactoryLaunched struct {
	Token       common.Address
	PoolId      [32]byte
	LaunchpadId [32]byte
	Creator     common.Address
	Quote       common.Address
	Supply      *big.Int
	TickSpacing *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterLaunched is a free log retrieval operation binding the contract event 0x0905ebd7b2bd18bf4fd23615c37069f5bb280f93c1ec6765c7c3d7054c41a2dc.
//
// Solidity: event Launched(address indexed token, bytes32 indexed poolId, bytes32 indexed launchpadId, address creator, address quote, uint256 supply, int24 tickSpacing)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryFilterer) FilterLaunched(opts *bind.FilterOpts, token []common.Address, poolId [][32]byte, launchpadId [][32]byte) (*MultiTenantLaunchpadFactoryLaunchedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var launchpadIdRule []interface{}
	for _, launchpadIdItem := range launchpadId {
		launchpadIdRule = append(launchpadIdRule, launchpadIdItem)
	}

	logs, sub, err := _MultiTenantLaunchpadFactory.contract.FilterLogs(opts, "Launched", tokenRule, poolIdRule, launchpadIdRule)
	if err != nil {
		return nil, err
	}
	return &MultiTenantLaunchpadFactoryLaunchedIterator{contract: _MultiTenantLaunchpadFactory.contract, event: "Launched", logs: logs, sub: sub}, nil
}

// WatchLaunched is a free log subscription operation binding the contract event 0x0905ebd7b2bd18bf4fd23615c37069f5bb280f93c1ec6765c7c3d7054c41a2dc.
//
// Solidity: event Launched(address indexed token, bytes32 indexed poolId, bytes32 indexed launchpadId, address creator, address quote, uint256 supply, int24 tickSpacing)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryFilterer) WatchLaunched(opts *bind.WatchOpts, sink chan<- *MultiTenantLaunchpadFactoryLaunched, token []common.Address, poolId [][32]byte, launchpadId [][32]byte) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var poolIdRule []interface{}
	for _, poolIdItem := range poolId {
		poolIdRule = append(poolIdRule, poolIdItem)
	}
	var launchpadIdRule []interface{}
	for _, launchpadIdItem := range launchpadId {
		launchpadIdRule = append(launchpadIdRule, launchpadIdItem)
	}

	logs, sub, err := _MultiTenantLaunchpadFactory.contract.WatchLogs(opts, "Launched", tokenRule, poolIdRule, launchpadIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultiTenantLaunchpadFactoryLaunched)
				if err := _MultiTenantLaunchpadFactory.contract.UnpackLog(event, "Launched", log); err != nil {
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

// ParseLaunched is a log parse operation binding the contract event 0x0905ebd7b2bd18bf4fd23615c37069f5bb280f93c1ec6765c7c3d7054c41a2dc.
//
// Solidity: event Launched(address indexed token, bytes32 indexed poolId, bytes32 indexed launchpadId, address creator, address quote, uint256 supply, int24 tickSpacing)
func (_MultiTenantLaunchpadFactory *MultiTenantLaunchpadFactoryFilterer) ParseLaunched(log types.Log) (*MultiTenantLaunchpadFactoryLaunched, error) {
	event := new(MultiTenantLaunchpadFactoryLaunched)
	if err := _MultiTenantLaunchpadFactory.contract.UnpackLog(event, "Launched", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
