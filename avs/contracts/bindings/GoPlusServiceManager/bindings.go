// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractGoPlusServiceManager

import (
	"errors"
	"math/big"
	"strings"

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
)

// ISignatureUtilsSignatureWithSaltAndExpiry is an auto generated low-level Go binding around an user-defined struct.
type ISignatureUtilsSignatureWithSaltAndExpiry struct {
	Signature []byte
	Salt      [32]byte
	Expiry    *big.Int
}

// GoPlusServiceManagerMetaData contains all meta data concerning the GoPlusServiceManager contract.
var GoPlusServiceManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_avsDirectory\",\"type\":\"address\",\"internalType\":\"contractIAVSDirectory\"},{\"name\":\"_registryCoordinator\",\"type\":\"address\",\"internalType\":\"contractIRegistryCoordinator\"},{\"name\":\"_stakeRegistry\",\"type\":\"address\",\"internalType\":\"contractIStakeRegistry\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"avsDirectory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deregisterOperatorFromAVS\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"gatewayAddr\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"gatewayURI\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorRestakedStrategies\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRestakeableStrategies\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_gatewayAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_gatewayURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_metadataURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerOperatorToAVS\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operatorSignature\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithSaltAndExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateAVSMetadataURI\",\"inputs\":[{\"name\":\"_metadataURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateGatewayAddress\",\"inputs\":[{\"name\":\"_gatewayAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateGatewayURI\",\"inputs\":[{\"name\":\"_gatewayURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"GatewayAddressUpdated\",\"inputs\":[{\"name\":\"oldAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GatewayURIUpdated\",\"inputs\":[{\"name\":\"oldURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"newURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false}]",
}

// GoPlusServiceManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use GoPlusServiceManagerMetaData.ABI instead.
var GoPlusServiceManagerABI = GoPlusServiceManagerMetaData.ABI

// GoPlusServiceManager is an auto generated Go binding around an Ethereum contract.
type GoPlusServiceManager struct {
	GoPlusServiceManagerCaller     // Read-only binding to the contract
	GoPlusServiceManagerTransactor // Write-only binding to the contract
	GoPlusServiceManagerFilterer   // Log filterer for contract events
}

// GoPlusServiceManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type GoPlusServiceManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GoPlusServiceManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GoPlusServiceManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GoPlusServiceManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GoPlusServiceManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GoPlusServiceManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GoPlusServiceManagerSession struct {
	Contract     *GoPlusServiceManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// GoPlusServiceManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GoPlusServiceManagerCallerSession struct {
	Contract *GoPlusServiceManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// GoPlusServiceManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GoPlusServiceManagerTransactorSession struct {
	Contract     *GoPlusServiceManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// GoPlusServiceManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type GoPlusServiceManagerRaw struct {
	Contract *GoPlusServiceManager // Generic contract binding to access the raw methods on
}

// GoPlusServiceManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GoPlusServiceManagerCallerRaw struct {
	Contract *GoPlusServiceManagerCaller // Generic read-only contract binding to access the raw methods on
}

// GoPlusServiceManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GoPlusServiceManagerTransactorRaw struct {
	Contract *GoPlusServiceManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGoPlusServiceManager creates a new instance of GoPlusServiceManager, bound to a specific deployed contract.
func NewGoPlusServiceManager(address common.Address, backend bind.ContractBackend) (*GoPlusServiceManager, error) {
	contract, err := bindGoPlusServiceManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GoPlusServiceManager{GoPlusServiceManagerCaller: GoPlusServiceManagerCaller{contract: contract}, GoPlusServiceManagerTransactor: GoPlusServiceManagerTransactor{contract: contract}, GoPlusServiceManagerFilterer: GoPlusServiceManagerFilterer{contract: contract}}, nil
}

// NewGoPlusServiceManagerCaller creates a new read-only instance of GoPlusServiceManager, bound to a specific deployed contract.
func NewGoPlusServiceManagerCaller(address common.Address, caller bind.ContractCaller) (*GoPlusServiceManagerCaller, error) {
	contract, err := bindGoPlusServiceManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GoPlusServiceManagerCaller{contract: contract}, nil
}

// NewGoPlusServiceManagerTransactor creates a new write-only instance of GoPlusServiceManager, bound to a specific deployed contract.
func NewGoPlusServiceManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*GoPlusServiceManagerTransactor, error) {
	contract, err := bindGoPlusServiceManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GoPlusServiceManagerTransactor{contract: contract}, nil
}

// NewGoPlusServiceManagerFilterer creates a new log filterer instance of GoPlusServiceManager, bound to a specific deployed contract.
func NewGoPlusServiceManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*GoPlusServiceManagerFilterer, error) {
	contract, err := bindGoPlusServiceManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GoPlusServiceManagerFilterer{contract: contract}, nil
}

// bindGoPlusServiceManager binds a generic wrapper to an already deployed contract.
func bindGoPlusServiceManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GoPlusServiceManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GoPlusServiceManager *GoPlusServiceManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GoPlusServiceManager.Contract.GoPlusServiceManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GoPlusServiceManager *GoPlusServiceManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.GoPlusServiceManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GoPlusServiceManager *GoPlusServiceManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.GoPlusServiceManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GoPlusServiceManager *GoPlusServiceManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GoPlusServiceManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GoPlusServiceManager *GoPlusServiceManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GoPlusServiceManager *GoPlusServiceManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.contract.Transact(opts, method, params...)
}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_GoPlusServiceManager *GoPlusServiceManagerCaller) AvsDirectory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GoPlusServiceManager.contract.Call(opts, &out, "avsDirectory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_GoPlusServiceManager *GoPlusServiceManagerSession) AvsDirectory() (common.Address, error) {
	return _GoPlusServiceManager.Contract.AvsDirectory(&_GoPlusServiceManager.CallOpts)
}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_GoPlusServiceManager *GoPlusServiceManagerCallerSession) AvsDirectory() (common.Address, error) {
	return _GoPlusServiceManager.Contract.AvsDirectory(&_GoPlusServiceManager.CallOpts)
}

// GatewayAddr is a free data retrieval call binding the contract method 0x3c46d619.
//
// Solidity: function gatewayAddr() view returns(address)
func (_GoPlusServiceManager *GoPlusServiceManagerCaller) GatewayAddr(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GoPlusServiceManager.contract.Call(opts, &out, "gatewayAddr")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GatewayAddr is a free data retrieval call binding the contract method 0x3c46d619.
//
// Solidity: function gatewayAddr() view returns(address)
func (_GoPlusServiceManager *GoPlusServiceManagerSession) GatewayAddr() (common.Address, error) {
	return _GoPlusServiceManager.Contract.GatewayAddr(&_GoPlusServiceManager.CallOpts)
}

// GatewayAddr is a free data retrieval call binding the contract method 0x3c46d619.
//
// Solidity: function gatewayAddr() view returns(address)
func (_GoPlusServiceManager *GoPlusServiceManagerCallerSession) GatewayAddr() (common.Address, error) {
	return _GoPlusServiceManager.Contract.GatewayAddr(&_GoPlusServiceManager.CallOpts)
}

// GatewayURI is a free data retrieval call binding the contract method 0x2f5199fb.
//
// Solidity: function gatewayURI() view returns(string)
func (_GoPlusServiceManager *GoPlusServiceManagerCaller) GatewayURI(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GoPlusServiceManager.contract.Call(opts, &out, "gatewayURI")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GatewayURI is a free data retrieval call binding the contract method 0x2f5199fb.
//
// Solidity: function gatewayURI() view returns(string)
func (_GoPlusServiceManager *GoPlusServiceManagerSession) GatewayURI() (string, error) {
	return _GoPlusServiceManager.Contract.GatewayURI(&_GoPlusServiceManager.CallOpts)
}

// GatewayURI is a free data retrieval call binding the contract method 0x2f5199fb.
//
// Solidity: function gatewayURI() view returns(string)
func (_GoPlusServiceManager *GoPlusServiceManagerCallerSession) GatewayURI() (string, error) {
	return _GoPlusServiceManager.Contract.GatewayURI(&_GoPlusServiceManager.CallOpts)
}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_GoPlusServiceManager *GoPlusServiceManagerCaller) GetOperatorRestakedStrategies(opts *bind.CallOpts, operator common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _GoPlusServiceManager.contract.Call(opts, &out, "getOperatorRestakedStrategies", operator)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_GoPlusServiceManager *GoPlusServiceManagerSession) GetOperatorRestakedStrategies(operator common.Address) ([]common.Address, error) {
	return _GoPlusServiceManager.Contract.GetOperatorRestakedStrategies(&_GoPlusServiceManager.CallOpts, operator)
}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_GoPlusServiceManager *GoPlusServiceManagerCallerSession) GetOperatorRestakedStrategies(operator common.Address) ([]common.Address, error) {
	return _GoPlusServiceManager.Contract.GetOperatorRestakedStrategies(&_GoPlusServiceManager.CallOpts, operator)
}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_GoPlusServiceManager *GoPlusServiceManagerCaller) GetRestakeableStrategies(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _GoPlusServiceManager.contract.Call(opts, &out, "getRestakeableStrategies")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_GoPlusServiceManager *GoPlusServiceManagerSession) GetRestakeableStrategies() ([]common.Address, error) {
	return _GoPlusServiceManager.Contract.GetRestakeableStrategies(&_GoPlusServiceManager.CallOpts)
}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_GoPlusServiceManager *GoPlusServiceManagerCallerSession) GetRestakeableStrategies() ([]common.Address, error) {
	return _GoPlusServiceManager.Contract.GetRestakeableStrategies(&_GoPlusServiceManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GoPlusServiceManager *GoPlusServiceManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GoPlusServiceManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GoPlusServiceManager *GoPlusServiceManagerSession) Owner() (common.Address, error) {
	return _GoPlusServiceManager.Contract.Owner(&_GoPlusServiceManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GoPlusServiceManager *GoPlusServiceManagerCallerSession) Owner() (common.Address, error) {
	return _GoPlusServiceManager.Contract.Owner(&_GoPlusServiceManager.CallOpts)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerTransactor) DeregisterOperatorFromAVS(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _GoPlusServiceManager.contract.Transact(opts, "deregisterOperatorFromAVS", operator)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerSession) DeregisterOperatorFromAVS(operator common.Address) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.DeregisterOperatorFromAVS(&_GoPlusServiceManager.TransactOpts, operator)
}

// DeregisterOperatorFromAVS is a paid mutator transaction binding the contract method 0xa364f4da.
//
// Solidity: function deregisterOperatorFromAVS(address operator) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerTransactorSession) DeregisterOperatorFromAVS(operator common.Address) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.DeregisterOperatorFromAVS(&_GoPlusServiceManager.TransactOpts, operator)
}

// Initialize is a paid mutator transaction binding the contract method 0x2016a0d2.
//
// Solidity: function initialize(address initialOwner, address _gatewayAddr, string _gatewayURI, string _metadataURI) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerTransactor) Initialize(opts *bind.TransactOpts, initialOwner common.Address, _gatewayAddr common.Address, _gatewayURI string, _metadataURI string) (*types.Transaction, error) {
	return _GoPlusServiceManager.contract.Transact(opts, "initialize", initialOwner, _gatewayAddr, _gatewayURI, _metadataURI)
}

// Initialize is a paid mutator transaction binding the contract method 0x2016a0d2.
//
// Solidity: function initialize(address initialOwner, address _gatewayAddr, string _gatewayURI, string _metadataURI) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerSession) Initialize(initialOwner common.Address, _gatewayAddr common.Address, _gatewayURI string, _metadataURI string) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.Initialize(&_GoPlusServiceManager.TransactOpts, initialOwner, _gatewayAddr, _gatewayURI, _metadataURI)
}

// Initialize is a paid mutator transaction binding the contract method 0x2016a0d2.
//
// Solidity: function initialize(address initialOwner, address _gatewayAddr, string _gatewayURI, string _metadataURI) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerTransactorSession) Initialize(initialOwner common.Address, _gatewayAddr common.Address, _gatewayURI string, _metadataURI string) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.Initialize(&_GoPlusServiceManager.TransactOpts, initialOwner, _gatewayAddr, _gatewayURI, _metadataURI)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerTransactor) RegisterOperatorToAVS(opts *bind.TransactOpts, operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _GoPlusServiceManager.contract.Transact(opts, "registerOperatorToAVS", operator, operatorSignature)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerSession) RegisterOperatorToAVS(operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.RegisterOperatorToAVS(&_GoPlusServiceManager.TransactOpts, operator, operatorSignature)
}

// RegisterOperatorToAVS is a paid mutator transaction binding the contract method 0x9926ee7d.
//
// Solidity: function registerOperatorToAVS(address operator, (bytes,bytes32,uint256) operatorSignature) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerTransactorSession) RegisterOperatorToAVS(operator common.Address, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.RegisterOperatorToAVS(&_GoPlusServiceManager.TransactOpts, operator, operatorSignature)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GoPlusServiceManager *GoPlusServiceManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GoPlusServiceManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GoPlusServiceManager *GoPlusServiceManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.RenounceOwnership(&_GoPlusServiceManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GoPlusServiceManager *GoPlusServiceManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.RenounceOwnership(&_GoPlusServiceManager.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _GoPlusServiceManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.TransferOwnership(&_GoPlusServiceManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.TransferOwnership(&_GoPlusServiceManager.TransactOpts, newOwner)
}

// UpdateAVSMetadataURI is a paid mutator transaction binding the contract method 0xa98fb355.
//
// Solidity: function updateAVSMetadataURI(string _metadataURI) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerTransactor) UpdateAVSMetadataURI(opts *bind.TransactOpts, _metadataURI string) (*types.Transaction, error) {
	return _GoPlusServiceManager.contract.Transact(opts, "updateAVSMetadataURI", _metadataURI)
}

// UpdateAVSMetadataURI is a paid mutator transaction binding the contract method 0xa98fb355.
//
// Solidity: function updateAVSMetadataURI(string _metadataURI) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerSession) UpdateAVSMetadataURI(_metadataURI string) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.UpdateAVSMetadataURI(&_GoPlusServiceManager.TransactOpts, _metadataURI)
}

// UpdateAVSMetadataURI is a paid mutator transaction binding the contract method 0xa98fb355.
//
// Solidity: function updateAVSMetadataURI(string _metadataURI) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerTransactorSession) UpdateAVSMetadataURI(_metadataURI string) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.UpdateAVSMetadataURI(&_GoPlusServiceManager.TransactOpts, _metadataURI)
}

// UpdateGatewayAddress is a paid mutator transaction binding the contract method 0xccc77599.
//
// Solidity: function updateGatewayAddress(address _gatewayAddr) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerTransactor) UpdateGatewayAddress(opts *bind.TransactOpts, _gatewayAddr common.Address) (*types.Transaction, error) {
	return _GoPlusServiceManager.contract.Transact(opts, "updateGatewayAddress", _gatewayAddr)
}

// UpdateGatewayAddress is a paid mutator transaction binding the contract method 0xccc77599.
//
// Solidity: function updateGatewayAddress(address _gatewayAddr) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerSession) UpdateGatewayAddress(_gatewayAddr common.Address) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.UpdateGatewayAddress(&_GoPlusServiceManager.TransactOpts, _gatewayAddr)
}

// UpdateGatewayAddress is a paid mutator transaction binding the contract method 0xccc77599.
//
// Solidity: function updateGatewayAddress(address _gatewayAddr) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerTransactorSession) UpdateGatewayAddress(_gatewayAddr common.Address) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.UpdateGatewayAddress(&_GoPlusServiceManager.TransactOpts, _gatewayAddr)
}

// UpdateGatewayURI is a paid mutator transaction binding the contract method 0x2bf005e3.
//
// Solidity: function updateGatewayURI(string _gatewayURI) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerTransactor) UpdateGatewayURI(opts *bind.TransactOpts, _gatewayURI string) (*types.Transaction, error) {
	return _GoPlusServiceManager.contract.Transact(opts, "updateGatewayURI", _gatewayURI)
}

// UpdateGatewayURI is a paid mutator transaction binding the contract method 0x2bf005e3.
//
// Solidity: function updateGatewayURI(string _gatewayURI) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerSession) UpdateGatewayURI(_gatewayURI string) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.UpdateGatewayURI(&_GoPlusServiceManager.TransactOpts, _gatewayURI)
}

// UpdateGatewayURI is a paid mutator transaction binding the contract method 0x2bf005e3.
//
// Solidity: function updateGatewayURI(string _gatewayURI) returns()
func (_GoPlusServiceManager *GoPlusServiceManagerTransactorSession) UpdateGatewayURI(_gatewayURI string) (*types.Transaction, error) {
	return _GoPlusServiceManager.Contract.UpdateGatewayURI(&_GoPlusServiceManager.TransactOpts, _gatewayURI)
}

// GoPlusServiceManagerGatewayAddressUpdatedIterator is returned from FilterGatewayAddressUpdated and is used to iterate over the raw logs and unpacked data for GatewayAddressUpdated events raised by the GoPlusServiceManager contract.
type GoPlusServiceManagerGatewayAddressUpdatedIterator struct {
	Event *GoPlusServiceManagerGatewayAddressUpdated // Event containing the contract specifics and raw log

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
func (it *GoPlusServiceManagerGatewayAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GoPlusServiceManagerGatewayAddressUpdated)
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
		it.Event = new(GoPlusServiceManagerGatewayAddressUpdated)
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
func (it *GoPlusServiceManagerGatewayAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GoPlusServiceManagerGatewayAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GoPlusServiceManagerGatewayAddressUpdated represents a GatewayAddressUpdated event raised by the GoPlusServiceManager contract.
type GoPlusServiceManagerGatewayAddressUpdated struct {
	OldAddr common.Address
	NewAddr common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterGatewayAddressUpdated is a free log retrieval operation binding the contract event 0xd4c520edf96d8835d69a539bded1b9d2b881f5e78ee3a66bb1d13e12013b5241.
//
// Solidity: event GatewayAddressUpdated(address indexed oldAddr, address indexed newAddr)
func (_GoPlusServiceManager *GoPlusServiceManagerFilterer) FilterGatewayAddressUpdated(opts *bind.FilterOpts, oldAddr []common.Address, newAddr []common.Address) (*GoPlusServiceManagerGatewayAddressUpdatedIterator, error) {

	var oldAddrRule []interface{}
	for _, oldAddrItem := range oldAddr {
		oldAddrRule = append(oldAddrRule, oldAddrItem)
	}
	var newAddrRule []interface{}
	for _, newAddrItem := range newAddr {
		newAddrRule = append(newAddrRule, newAddrItem)
	}

	logs, sub, err := _GoPlusServiceManager.contract.FilterLogs(opts, "GatewayAddressUpdated", oldAddrRule, newAddrRule)
	if err != nil {
		return nil, err
	}
	return &GoPlusServiceManagerGatewayAddressUpdatedIterator{contract: _GoPlusServiceManager.contract, event: "GatewayAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchGatewayAddressUpdated is a free log subscription operation binding the contract event 0xd4c520edf96d8835d69a539bded1b9d2b881f5e78ee3a66bb1d13e12013b5241.
//
// Solidity: event GatewayAddressUpdated(address indexed oldAddr, address indexed newAddr)
func (_GoPlusServiceManager *GoPlusServiceManagerFilterer) WatchGatewayAddressUpdated(opts *bind.WatchOpts, sink chan<- *GoPlusServiceManagerGatewayAddressUpdated, oldAddr []common.Address, newAddr []common.Address) (event.Subscription, error) {

	var oldAddrRule []interface{}
	for _, oldAddrItem := range oldAddr {
		oldAddrRule = append(oldAddrRule, oldAddrItem)
	}
	var newAddrRule []interface{}
	for _, newAddrItem := range newAddr {
		newAddrRule = append(newAddrRule, newAddrItem)
	}

	logs, sub, err := _GoPlusServiceManager.contract.WatchLogs(opts, "GatewayAddressUpdated", oldAddrRule, newAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GoPlusServiceManagerGatewayAddressUpdated)
				if err := _GoPlusServiceManager.contract.UnpackLog(event, "GatewayAddressUpdated", log); err != nil {
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

// ParseGatewayAddressUpdated is a log parse operation binding the contract event 0xd4c520edf96d8835d69a539bded1b9d2b881f5e78ee3a66bb1d13e12013b5241.
//
// Solidity: event GatewayAddressUpdated(address indexed oldAddr, address indexed newAddr)
func (_GoPlusServiceManager *GoPlusServiceManagerFilterer) ParseGatewayAddressUpdated(log types.Log) (*GoPlusServiceManagerGatewayAddressUpdated, error) {
	event := new(GoPlusServiceManagerGatewayAddressUpdated)
	if err := _GoPlusServiceManager.contract.UnpackLog(event, "GatewayAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GoPlusServiceManagerGatewayURIUpdatedIterator is returned from FilterGatewayURIUpdated and is used to iterate over the raw logs and unpacked data for GatewayURIUpdated events raised by the GoPlusServiceManager contract.
type GoPlusServiceManagerGatewayURIUpdatedIterator struct {
	Event *GoPlusServiceManagerGatewayURIUpdated // Event containing the contract specifics and raw log

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
func (it *GoPlusServiceManagerGatewayURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GoPlusServiceManagerGatewayURIUpdated)
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
		it.Event = new(GoPlusServiceManagerGatewayURIUpdated)
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
func (it *GoPlusServiceManagerGatewayURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GoPlusServiceManagerGatewayURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GoPlusServiceManagerGatewayURIUpdated represents a GatewayURIUpdated event raised by the GoPlusServiceManager contract.
type GoPlusServiceManagerGatewayURIUpdated struct {
	OldURI string
	NewURI string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterGatewayURIUpdated is a free log retrieval operation binding the contract event 0xe901d1036a5f8141a1503434545c3c097f1d7e26e1e9e74e7e43e98f1d75e892.
//
// Solidity: event GatewayURIUpdated(string oldURI, string newURI)
func (_GoPlusServiceManager *GoPlusServiceManagerFilterer) FilterGatewayURIUpdated(opts *bind.FilterOpts) (*GoPlusServiceManagerGatewayURIUpdatedIterator, error) {

	logs, sub, err := _GoPlusServiceManager.contract.FilterLogs(opts, "GatewayURIUpdated")
	if err != nil {
		return nil, err
	}
	return &GoPlusServiceManagerGatewayURIUpdatedIterator{contract: _GoPlusServiceManager.contract, event: "GatewayURIUpdated", logs: logs, sub: sub}, nil
}

// WatchGatewayURIUpdated is a free log subscription operation binding the contract event 0xe901d1036a5f8141a1503434545c3c097f1d7e26e1e9e74e7e43e98f1d75e892.
//
// Solidity: event GatewayURIUpdated(string oldURI, string newURI)
func (_GoPlusServiceManager *GoPlusServiceManagerFilterer) WatchGatewayURIUpdated(opts *bind.WatchOpts, sink chan<- *GoPlusServiceManagerGatewayURIUpdated) (event.Subscription, error) {

	logs, sub, err := _GoPlusServiceManager.contract.WatchLogs(opts, "GatewayURIUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GoPlusServiceManagerGatewayURIUpdated)
				if err := _GoPlusServiceManager.contract.UnpackLog(event, "GatewayURIUpdated", log); err != nil {
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

// ParseGatewayURIUpdated is a log parse operation binding the contract event 0xe901d1036a5f8141a1503434545c3c097f1d7e26e1e9e74e7e43e98f1d75e892.
//
// Solidity: event GatewayURIUpdated(string oldURI, string newURI)
func (_GoPlusServiceManager *GoPlusServiceManagerFilterer) ParseGatewayURIUpdated(log types.Log) (*GoPlusServiceManagerGatewayURIUpdated, error) {
	event := new(GoPlusServiceManagerGatewayURIUpdated)
	if err := _GoPlusServiceManager.contract.UnpackLog(event, "GatewayURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GoPlusServiceManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the GoPlusServiceManager contract.
type GoPlusServiceManagerInitializedIterator struct {
	Event *GoPlusServiceManagerInitialized // Event containing the contract specifics and raw log

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
func (it *GoPlusServiceManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GoPlusServiceManagerInitialized)
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
		it.Event = new(GoPlusServiceManagerInitialized)
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
func (it *GoPlusServiceManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GoPlusServiceManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GoPlusServiceManagerInitialized represents a Initialized event raised by the GoPlusServiceManager contract.
type GoPlusServiceManagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GoPlusServiceManager *GoPlusServiceManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*GoPlusServiceManagerInitializedIterator, error) {

	logs, sub, err := _GoPlusServiceManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &GoPlusServiceManagerInitializedIterator{contract: _GoPlusServiceManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GoPlusServiceManager *GoPlusServiceManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *GoPlusServiceManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _GoPlusServiceManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GoPlusServiceManagerInitialized)
				if err := _GoPlusServiceManager.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GoPlusServiceManager *GoPlusServiceManagerFilterer) ParseInitialized(log types.Log) (*GoPlusServiceManagerInitialized, error) {
	event := new(GoPlusServiceManagerInitialized)
	if err := _GoPlusServiceManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GoPlusServiceManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the GoPlusServiceManager contract.
type GoPlusServiceManagerOwnershipTransferredIterator struct {
	Event *GoPlusServiceManagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *GoPlusServiceManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GoPlusServiceManagerOwnershipTransferred)
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
		it.Event = new(GoPlusServiceManagerOwnershipTransferred)
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
func (it *GoPlusServiceManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GoPlusServiceManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GoPlusServiceManagerOwnershipTransferred represents a OwnershipTransferred event raised by the GoPlusServiceManager contract.
type GoPlusServiceManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GoPlusServiceManager *GoPlusServiceManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*GoPlusServiceManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GoPlusServiceManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &GoPlusServiceManagerOwnershipTransferredIterator{contract: _GoPlusServiceManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GoPlusServiceManager *GoPlusServiceManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *GoPlusServiceManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GoPlusServiceManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GoPlusServiceManagerOwnershipTransferred)
				if err := _GoPlusServiceManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GoPlusServiceManager *GoPlusServiceManagerFilterer) ParseOwnershipTransferred(log types.Log) (*GoPlusServiceManagerOwnershipTransferred, error) {
	event := new(GoPlusServiceManagerOwnershipTransferred)
	if err := _GoPlusServiceManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
