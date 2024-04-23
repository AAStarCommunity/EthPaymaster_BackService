// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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

// OracleHelperOracleHelperConfig is an auto generated low-level Go binding around an user-defined struct.
type OracleHelperOracleHelperConfig struct {
	CacheTimeToLive      *big.Int
	MaxOracleRoundAge    *big.Int
	TokenOracle          common.Address
	NativeOracle         common.Address
	TokenToNativeOracle  bool
	TokenOracleReverse   bool
	NativeOracleReverse  bool
	PriceUpdateThreshold *big.Int
}

// PackedUserOperation is an auto generated low-level Go binding around an user-defined struct.
type PackedUserOperation struct {
	Sender             common.Address
	Nonce              *big.Int
	InitCode           []byte
	CallData           []byte
	AccountGasLimits   [32]byte
	PreVerificationGas *big.Int
	GasFees            [32]byte
	PaymasterAndData   []byte
	Signature          []byte
}

// TokenPaymasterTokenPaymasterConfig is an auto generated low-level Go binding around an user-defined struct.
type TokenPaymasterTokenPaymasterConfig struct {
	PriceMarkup          *big.Int
	MinEntryPointBalance *big.Int
	RefundPostopCost     *big.Int
	PriceMaxAge          *big.Int
}

// UniswapHelperUniswapHelperConfig is an auto generated low-level Go binding around an user-defined struct.
type UniswapHelperUniswapHelperConfig struct {
	MinSwapAmount  *big.Int
	UniswapPoolFee *big.Int
	Slippage       uint8
}

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20Metadata\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"contractIEntryPoint\",\"name\":\"_entryPoint\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"_wrappedNative\",\"type\":\"address\"},{\"internalType\":\"contractISwapRouter\",\"name\":\"_uniswap\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"priceMarkup\",\"type\":\"uint256\"},{\"internalType\":\"uint128\",\"name\":\"minEntryPointBalance\",\"type\":\"uint128\"},{\"internalType\":\"uint48\",\"name\":\"refundPostopCost\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"priceMaxAge\",\"type\":\"uint48\"}],\"internalType\":\"structTokenPaymaster.TokenPaymasterConfig\",\"name\":\"_tokenPaymasterConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint48\",\"name\":\"cacheTimeToLive\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"maxOracleRoundAge\",\"type\":\"uint48\"},{\"internalType\":\"contractIOracle\",\"name\":\"tokenOracle\",\"type\":\"address\"},{\"internalType\":\"contractIOracle\",\"name\":\"nativeOracle\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"tokenToNativeOracle\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"tokenOracleReverse\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"nativeOracleReverse\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"priceUpdateThreshold\",\"type\":\"uint256\"}],\"internalType\":\"structOracleHelper.OracleHelperConfig\",\"name\":\"_oracleHelperConfig\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"minSwapAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint24\",\"name\":\"uniswapPoolFee\",\"type\":\"uint24\"},{\"internalType\":\"uint8\",\"name\":\"slippage\",\"type\":\"uint8\"}],\"internalType\":\"structUniswapHelper.UniswapHelperConfig\",\"name\":\"_uniswapHelperConfig\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"priceMarkup\",\"type\":\"uint256\"},{\"internalType\":\"uint128\",\"name\":\"minEntryPointBalance\",\"type\":\"uint128\"},{\"internalType\":\"uint48\",\"name\":\"refundPostopCost\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"priceMaxAge\",\"type\":\"uint48\"}],\"indexed\":false,\"internalType\":\"structTokenPaymaster.TokenPaymasterConfig\",\"name\":\"tokenPaymasterConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Received\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"currentPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"previousPrice\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"cachedPriceTimestamp\",\"type\":\"uint256\"}],\"name\":\"TokenPriceUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"}],\"name\":\"UniswapReverted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"actualTokenCharge\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"actualGasCost\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"actualTokenPriceWithMarkup\",\"type\":\"uint256\"}],\"name\":\"UserOperationSponsored\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"unstakeDelaySec\",\"type\":\"uint32\"}],\"name\":\"addStake\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cachedPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cachedPriceTimestamp\",\"outputs\":[{\"internalType\":\"uint48\",\"name\":\"\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"entryPoint\",\"outputs\":[{\"internalType\":\"contractIEntryPoint\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDeposit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumIPaymaster.PostOpMode\",\"name\":\"mode\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"context\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"actualGasCost\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualUserOpFeePerGas\",\"type\":\"uint256\"}],\"name\":\"postOp\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"priceMarkup\",\"type\":\"uint256\"},{\"internalType\":\"uint128\",\"name\":\"minEntryPointBalance\",\"type\":\"uint128\"},{\"internalType\":\"uint48\",\"name\":\"refundPostopCost\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"priceMaxAge\",\"type\":\"uint48\"}],\"internalType\":\"structTokenPaymaster.TokenPaymasterConfig\",\"name\":\"_tokenPaymasterConfig\",\"type\":\"tuple\"}],\"name\":\"setTokenPaymasterConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"minSwapAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint24\",\"name\":\"uniswapPoolFee\",\"type\":\"uint24\"},{\"internalType\":\"uint8\",\"name\":\"slippage\",\"type\":\"uint8\"}],\"internalType\":\"structUniswapHelper.UniswapHelperConfig\",\"name\":\"_uniswapHelperConfig\",\"type\":\"tuple\"}],\"name\":\"setUniswapConfiguration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenPaymasterConfig\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"priceMarkup\",\"type\":\"uint256\"},{\"internalType\":\"uint128\",\"name\":\"minEntryPointBalance\",\"type\":\"uint128\"},{\"internalType\":\"uint48\",\"name\":\"refundPostopCost\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"priceMaxAge\",\"type\":\"uint48\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"tokenToWei\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"uniswap\",\"outputs\":[{\"internalType\":\"contractISwapRouter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unlockStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"force\",\"type\":\"bool\"}],\"name\":\"updateCachedPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"initCode\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"accountGasLimits\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"preVerificationGas\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"gasFees\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"paymasterAndData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structPackedUserOperation\",\"name\":\"userOp\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"userOpHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maxCost\",\"type\":\"uint256\"}],\"name\":\"validatePaymasterUserOp\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"context\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"validationData\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"weiToToken\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawEth\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"withdrawAddress\",\"type\":\"address\"}],\"name\":\"withdrawStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"withdrawAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"wrappedNative\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// ContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMetaData.ABI instead.
var ContractABI = ContractMetaData.ABI

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// CachedPrice is a free data retrieval call binding the contract method 0xf60fdcb3.
//
// Solidity: function cachedPrice() view returns(uint256)
func (_Contract *ContractCaller) CachedPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "cachedPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CachedPrice is a free data retrieval call binding the contract method 0xf60fdcb3.
//
// Solidity: function cachedPrice() view returns(uint256)
func (_Contract *ContractSession) CachedPrice() (*big.Int, error) {
	return _Contract.Contract.CachedPrice(&_Contract.CallOpts)
}

// CachedPrice is a free data retrieval call binding the contract method 0xf60fdcb3.
//
// Solidity: function cachedPrice() view returns(uint256)
func (_Contract *ContractCallerSession) CachedPrice() (*big.Int, error) {
	return _Contract.Contract.CachedPrice(&_Contract.CallOpts)
}

// CachedPriceTimestamp is a free data retrieval call binding the contract method 0xe1d8153c.
//
// Solidity: function cachedPriceTimestamp() view returns(uint48)
func (_Contract *ContractCaller) CachedPriceTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "cachedPriceTimestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CachedPriceTimestamp is a free data retrieval call binding the contract method 0xe1d8153c.
//
// Solidity: function cachedPriceTimestamp() view returns(uint48)
func (_Contract *ContractSession) CachedPriceTimestamp() (*big.Int, error) {
	return _Contract.Contract.CachedPriceTimestamp(&_Contract.CallOpts)
}

// CachedPriceTimestamp is a free data retrieval call binding the contract method 0xe1d8153c.
//
// Solidity: function cachedPriceTimestamp() view returns(uint48)
func (_Contract *ContractCallerSession) CachedPriceTimestamp() (*big.Int, error) {
	return _Contract.Contract.CachedPriceTimestamp(&_Contract.CallOpts)
}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_Contract *ContractCaller) EntryPoint(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "entryPoint")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_Contract *ContractSession) EntryPoint() (common.Address, error) {
	return _Contract.Contract.EntryPoint(&_Contract.CallOpts)
}

// EntryPoint is a free data retrieval call binding the contract method 0xb0d691fe.
//
// Solidity: function entryPoint() view returns(address)
func (_Contract *ContractCallerSession) EntryPoint() (common.Address, error) {
	return _Contract.Contract.EntryPoint(&_Contract.CallOpts)
}

// GetDeposit is a free data retrieval call binding the contract method 0xc399ec88.
//
// Solidity: function getDeposit() view returns(uint256)
func (_Contract *ContractCaller) GetDeposit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getDeposit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDeposit is a free data retrieval call binding the contract method 0xc399ec88.
//
// Solidity: function getDeposit() view returns(uint256)
func (_Contract *ContractSession) GetDeposit() (*big.Int, error) {
	return _Contract.Contract.GetDeposit(&_Contract.CallOpts)
}

// GetDeposit is a free data retrieval call binding the contract method 0xc399ec88.
//
// Solidity: function getDeposit() view returns(uint256)
func (_Contract *ContractCallerSession) GetDeposit() (*big.Int, error) {
	return _Contract.Contract.GetDeposit(&_Contract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractSession) Owner() (common.Address, error) {
	return _Contract.Contract.Owner(&_Contract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractCallerSession) Owner() (common.Address, error) {
	return _Contract.Contract.Owner(&_Contract.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Contract *ContractCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Contract *ContractSession) Token() (common.Address, error) {
	return _Contract.Contract.Token(&_Contract.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Contract *ContractCallerSession) Token() (common.Address, error) {
	return _Contract.Contract.Token(&_Contract.CallOpts)
}

// TokenPaymasterConfig is a free data retrieval call binding the contract method 0xcb721cfd.
//
// Solidity: function tokenPaymasterConfig() view returns(uint256 priceMarkup, uint128 minEntryPointBalance, uint48 refundPostopCost, uint48 priceMaxAge)
func (_Contract *ContractCaller) TokenPaymasterConfig(opts *bind.CallOpts) (struct {
	PriceMarkup          *big.Int
	MinEntryPointBalance *big.Int
	RefundPostopCost     *big.Int
	PriceMaxAge          *big.Int
}, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "tokenPaymasterConfig")

	outstruct := new(struct {
		PriceMarkup          *big.Int
		MinEntryPointBalance *big.Int
		RefundPostopCost     *big.Int
		PriceMaxAge          *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.PriceMarkup = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.MinEntryPointBalance = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.RefundPostopCost = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.PriceMaxAge = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// TokenPaymasterConfig is a free data retrieval call binding the contract method 0xcb721cfd.
//
// Solidity: function tokenPaymasterConfig() view returns(uint256 priceMarkup, uint128 minEntryPointBalance, uint48 refundPostopCost, uint48 priceMaxAge)
func (_Contract *ContractSession) TokenPaymasterConfig() (struct {
	PriceMarkup          *big.Int
	MinEntryPointBalance *big.Int
	RefundPostopCost     *big.Int
	PriceMaxAge          *big.Int
}, error) {
	return _Contract.Contract.TokenPaymasterConfig(&_Contract.CallOpts)
}

// TokenPaymasterConfig is a free data retrieval call binding the contract method 0xcb721cfd.
//
// Solidity: function tokenPaymasterConfig() view returns(uint256 priceMarkup, uint128 minEntryPointBalance, uint48 refundPostopCost, uint48 priceMaxAge)
func (_Contract *ContractCallerSession) TokenPaymasterConfig() (struct {
	PriceMarkup          *big.Int
	MinEntryPointBalance *big.Int
	RefundPostopCost     *big.Int
	PriceMaxAge          *big.Int
}, error) {
	return _Contract.Contract.TokenPaymasterConfig(&_Contract.CallOpts)
}

// TokenToWei is a free data retrieval call binding the contract method 0xd7a23b3c.
//
// Solidity: function tokenToWei(uint256 amount, uint256 price) pure returns(uint256)
func (_Contract *ContractCaller) TokenToWei(opts *bind.CallOpts, amount *big.Int, price *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "tokenToWei", amount, price)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenToWei is a free data retrieval call binding the contract method 0xd7a23b3c.
//
// Solidity: function tokenToWei(uint256 amount, uint256 price) pure returns(uint256)
func (_Contract *ContractSession) TokenToWei(amount *big.Int, price *big.Int) (*big.Int, error) {
	return _Contract.Contract.TokenToWei(&_Contract.CallOpts, amount, price)
}

// TokenToWei is a free data retrieval call binding the contract method 0xd7a23b3c.
//
// Solidity: function tokenToWei(uint256 amount, uint256 price) pure returns(uint256)
func (_Contract *ContractCallerSession) TokenToWei(amount *big.Int, price *big.Int) (*big.Int, error) {
	return _Contract.Contract.TokenToWei(&_Contract.CallOpts, amount, price)
}

// Uniswap is a free data retrieval call binding the contract method 0x2681f7e4.
//
// Solidity: function uniswap() view returns(address)
func (_Contract *ContractCaller) Uniswap(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "uniswap")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Uniswap is a free data retrieval call binding the contract method 0x2681f7e4.
//
// Solidity: function uniswap() view returns(address)
func (_Contract *ContractSession) Uniswap() (common.Address, error) {
	return _Contract.Contract.Uniswap(&_Contract.CallOpts)
}

// Uniswap is a free data retrieval call binding the contract method 0x2681f7e4.
//
// Solidity: function uniswap() view returns(address)
func (_Contract *ContractCallerSession) Uniswap() (common.Address, error) {
	return _Contract.Contract.Uniswap(&_Contract.CallOpts)
}

// WeiToToken is a free data retrieval call binding the contract method 0x7c986aac.
//
// Solidity: function weiToToken(uint256 amount, uint256 price) pure returns(uint256)
func (_Contract *ContractCaller) WeiToToken(opts *bind.CallOpts, amount *big.Int, price *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "weiToToken", amount, price)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WeiToToken is a free data retrieval call binding the contract method 0x7c986aac.
//
// Solidity: function weiToToken(uint256 amount, uint256 price) pure returns(uint256)
func (_Contract *ContractSession) WeiToToken(amount *big.Int, price *big.Int) (*big.Int, error) {
	return _Contract.Contract.WeiToToken(&_Contract.CallOpts, amount, price)
}

// WeiToToken is a free data retrieval call binding the contract method 0x7c986aac.
//
// Solidity: function weiToToken(uint256 amount, uint256 price) pure returns(uint256)
func (_Contract *ContractCallerSession) WeiToToken(amount *big.Int, price *big.Int) (*big.Int, error) {
	return _Contract.Contract.WeiToToken(&_Contract.CallOpts, amount, price)
}

// WrappedNative is a free data retrieval call binding the contract method 0xeb6d3a11.
//
// Solidity: function wrappedNative() view returns(address)
func (_Contract *ContractCaller) WrappedNative(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "wrappedNative")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WrappedNative is a free data retrieval call binding the contract method 0xeb6d3a11.
//
// Solidity: function wrappedNative() view returns(address)
func (_Contract *ContractSession) WrappedNative() (common.Address, error) {
	return _Contract.Contract.WrappedNative(&_Contract.CallOpts)
}

// WrappedNative is a free data retrieval call binding the contract method 0xeb6d3a11.
//
// Solidity: function wrappedNative() view returns(address)
func (_Contract *ContractCallerSession) WrappedNative() (common.Address, error) {
	return _Contract.Contract.WrappedNative(&_Contract.CallOpts)
}

// AddStake is a paid mutator transaction binding the contract method 0x0396cb60.
//
// Solidity: function addStake(uint32 unstakeDelaySec) payable returns()
func (_Contract *ContractTransactor) AddStake(opts *bind.TransactOpts, unstakeDelaySec uint32) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "addStake", unstakeDelaySec)
}

// AddStake is a paid mutator transaction binding the contract method 0x0396cb60.
//
// Solidity: function addStake(uint32 unstakeDelaySec) payable returns()
func (_Contract *ContractSession) AddStake(unstakeDelaySec uint32) (*types.Transaction, error) {
	return _Contract.Contract.AddStake(&_Contract.TransactOpts, unstakeDelaySec)
}

// AddStake is a paid mutator transaction binding the contract method 0x0396cb60.
//
// Solidity: function addStake(uint32 unstakeDelaySec) payable returns()
func (_Contract *ContractTransactorSession) AddStake(unstakeDelaySec uint32) (*types.Transaction, error) {
	return _Contract.Contract.AddStake(&_Contract.TransactOpts, unstakeDelaySec)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_Contract *ContractTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_Contract *ContractSession) Deposit() (*types.Transaction, error) {
	return _Contract.Contract.Deposit(&_Contract.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_Contract *ContractTransactorSession) Deposit() (*types.Transaction, error) {
	return _Contract.Contract.Deposit(&_Contract.TransactOpts)
}

// PostOp is a paid mutator transaction binding the contract method 0x7c627b21.
//
// Solidity: function postOp(uint8 mode, bytes context, uint256 actualGasCost, uint256 actualUserOpFeePerGas) returns()
func (_Contract *ContractTransactor) PostOp(opts *bind.TransactOpts, mode uint8, context []byte, actualGasCost *big.Int, actualUserOpFeePerGas *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "postOp", mode, context, actualGasCost, actualUserOpFeePerGas)
}

// PostOp is a paid mutator transaction binding the contract method 0x7c627b21.
//
// Solidity: function postOp(uint8 mode, bytes context, uint256 actualGasCost, uint256 actualUserOpFeePerGas) returns()
func (_Contract *ContractSession) PostOp(mode uint8, context []byte, actualGasCost *big.Int, actualUserOpFeePerGas *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PostOp(&_Contract.TransactOpts, mode, context, actualGasCost, actualUserOpFeePerGas)
}

// PostOp is a paid mutator transaction binding the contract method 0x7c627b21.
//
// Solidity: function postOp(uint8 mode, bytes context, uint256 actualGasCost, uint256 actualUserOpFeePerGas) returns()
func (_Contract *ContractTransactorSession) PostOp(mode uint8, context []byte, actualGasCost *big.Int, actualUserOpFeePerGas *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PostOp(&_Contract.TransactOpts, mode, context, actualGasCost, actualUserOpFeePerGas)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Contract *ContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Contract *ContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _Contract.Contract.RenounceOwnership(&_Contract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Contract *ContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Contract.Contract.RenounceOwnership(&_Contract.TransactOpts)
}

// SetTokenPaymasterConfig is a paid mutator transaction binding the contract method 0xf14d64ed.
//
// Solidity: function setTokenPaymasterConfig((uint256,uint128,uint48,uint48) _tokenPaymasterConfig) returns()
func (_Contract *ContractTransactor) SetTokenPaymasterConfig(opts *bind.TransactOpts, _tokenPaymasterConfig TokenPaymasterTokenPaymasterConfig) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setTokenPaymasterConfig", _tokenPaymasterConfig)
}

// SetTokenPaymasterConfig is a paid mutator transaction binding the contract method 0xf14d64ed.
//
// Solidity: function setTokenPaymasterConfig((uint256,uint128,uint48,uint48) _tokenPaymasterConfig) returns()
func (_Contract *ContractSession) SetTokenPaymasterConfig(_tokenPaymasterConfig TokenPaymasterTokenPaymasterConfig) (*types.Transaction, error) {
	return _Contract.Contract.SetTokenPaymasterConfig(&_Contract.TransactOpts, _tokenPaymasterConfig)
}

// SetTokenPaymasterConfig is a paid mutator transaction binding the contract method 0xf14d64ed.
//
// Solidity: function setTokenPaymasterConfig((uint256,uint128,uint48,uint48) _tokenPaymasterConfig) returns()
func (_Contract *ContractTransactorSession) SetTokenPaymasterConfig(_tokenPaymasterConfig TokenPaymasterTokenPaymasterConfig) (*types.Transaction, error) {
	return _Contract.Contract.SetTokenPaymasterConfig(&_Contract.TransactOpts, _tokenPaymasterConfig)
}

// SetUniswapConfiguration is a paid mutator transaction binding the contract method 0xa0840fa7.
//
// Solidity: function setUniswapConfiguration((uint256,uint24,uint8) _uniswapHelperConfig) returns()
func (_Contract *ContractTransactor) SetUniswapConfiguration(opts *bind.TransactOpts, _uniswapHelperConfig UniswapHelperUniswapHelperConfig) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setUniswapConfiguration", _uniswapHelperConfig)
}

// SetUniswapConfiguration is a paid mutator transaction binding the contract method 0xa0840fa7.
//
// Solidity: function setUniswapConfiguration((uint256,uint24,uint8) _uniswapHelperConfig) returns()
func (_Contract *ContractSession) SetUniswapConfiguration(_uniswapHelperConfig UniswapHelperUniswapHelperConfig) (*types.Transaction, error) {
	return _Contract.Contract.SetUniswapConfiguration(&_Contract.TransactOpts, _uniswapHelperConfig)
}

// SetUniswapConfiguration is a paid mutator transaction binding the contract method 0xa0840fa7.
//
// Solidity: function setUniswapConfiguration((uint256,uint24,uint8) _uniswapHelperConfig) returns()
func (_Contract *ContractTransactorSession) SetUniswapConfiguration(_uniswapHelperConfig UniswapHelperUniswapHelperConfig) (*types.Transaction, error) {
	return _Contract.Contract.SetUniswapConfiguration(&_Contract.TransactOpts, _uniswapHelperConfig)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.TransferOwnership(&_Contract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.TransferOwnership(&_Contract.TransactOpts, newOwner)
}

// UnlockStake is a paid mutator transaction binding the contract method 0xbb9fe6bf.
//
// Solidity: function unlockStake() returns()
func (_Contract *ContractTransactor) UnlockStake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "unlockStake")
}

// UnlockStake is a paid mutator transaction binding the contract method 0xbb9fe6bf.
//
// Solidity: function unlockStake() returns()
func (_Contract *ContractSession) UnlockStake() (*types.Transaction, error) {
	return _Contract.Contract.UnlockStake(&_Contract.TransactOpts)
}

// UnlockStake is a paid mutator transaction binding the contract method 0xbb9fe6bf.
//
// Solidity: function unlockStake() returns()
func (_Contract *ContractTransactorSession) UnlockStake() (*types.Transaction, error) {
	return _Contract.Contract.UnlockStake(&_Contract.TransactOpts)
}

// UpdateCachedPrice is a paid mutator transaction binding the contract method 0x3ba9290f.
//
// Solidity: function updateCachedPrice(bool force) returns(uint256)
func (_Contract *ContractTransactor) UpdateCachedPrice(opts *bind.TransactOpts, force bool) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "updateCachedPrice", force)
}

// UpdateCachedPrice is a paid mutator transaction binding the contract method 0x3ba9290f.
//
// Solidity: function updateCachedPrice(bool force) returns(uint256)
func (_Contract *ContractSession) UpdateCachedPrice(force bool) (*types.Transaction, error) {
	return _Contract.Contract.UpdateCachedPrice(&_Contract.TransactOpts, force)
}

// UpdateCachedPrice is a paid mutator transaction binding the contract method 0x3ba9290f.
//
// Solidity: function updateCachedPrice(bool force) returns(uint256)
func (_Contract *ContractTransactorSession) UpdateCachedPrice(force bool) (*types.Transaction, error) {
	return _Contract.Contract.UpdateCachedPrice(&_Contract.TransactOpts, force)
}

// ValidatePaymasterUserOp is a paid mutator transaction binding the contract method 0x52b7512c.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 maxCost) returns(bytes context, uint256 validationData)
func (_Contract *ContractTransactor) ValidatePaymasterUserOp(opts *bind.TransactOpts, userOp PackedUserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "validatePaymasterUserOp", userOp, userOpHash, maxCost)
}

// ValidatePaymasterUserOp is a paid mutator transaction binding the contract method 0x52b7512c.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 maxCost) returns(bytes context, uint256 validationData)
func (_Contract *ContractSession) ValidatePaymasterUserOp(userOp PackedUserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.ValidatePaymasterUserOp(&_Contract.TransactOpts, userOp, userOpHash, maxCost)
}

// ValidatePaymasterUserOp is a paid mutator transaction binding the contract method 0x52b7512c.
//
// Solidity: function validatePaymasterUserOp((address,uint256,bytes,bytes,bytes32,uint256,bytes32,bytes,bytes) userOp, bytes32 userOpHash, uint256 maxCost) returns(bytes context, uint256 validationData)
func (_Contract *ContractTransactorSession) ValidatePaymasterUserOp(userOp PackedUserOperation, userOpHash [32]byte, maxCost *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.ValidatePaymasterUserOp(&_Contract.TransactOpts, userOp, userOpHash, maxCost)
}

// WithdrawEth is a paid mutator transaction binding the contract method 0x1b9a91a4.
//
// Solidity: function withdrawEth(address recipient, uint256 amount) returns()
func (_Contract *ContractTransactor) WithdrawEth(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "withdrawEth", recipient, amount)
}

// WithdrawEth is a paid mutator transaction binding the contract method 0x1b9a91a4.
//
// Solidity: function withdrawEth(address recipient, uint256 amount) returns()
func (_Contract *ContractSession) WithdrawEth(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.WithdrawEth(&_Contract.TransactOpts, recipient, amount)
}

// WithdrawEth is a paid mutator transaction binding the contract method 0x1b9a91a4.
//
// Solidity: function withdrawEth(address recipient, uint256 amount) returns()
func (_Contract *ContractTransactorSession) WithdrawEth(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.WithdrawEth(&_Contract.TransactOpts, recipient, amount)
}

// WithdrawStake is a paid mutator transaction binding the contract method 0xc23a5cea.
//
// Solidity: function withdrawStake(address withdrawAddress) returns()
func (_Contract *ContractTransactor) WithdrawStake(opts *bind.TransactOpts, withdrawAddress common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "withdrawStake", withdrawAddress)
}

// WithdrawStake is a paid mutator transaction binding the contract method 0xc23a5cea.
//
// Solidity: function withdrawStake(address withdrawAddress) returns()
func (_Contract *ContractSession) WithdrawStake(withdrawAddress common.Address) (*types.Transaction, error) {
	return _Contract.Contract.WithdrawStake(&_Contract.TransactOpts, withdrawAddress)
}

// WithdrawStake is a paid mutator transaction binding the contract method 0xc23a5cea.
//
// Solidity: function withdrawStake(address withdrawAddress) returns()
func (_Contract *ContractTransactorSession) WithdrawStake(withdrawAddress common.Address) (*types.Transaction, error) {
	return _Contract.Contract.WithdrawStake(&_Contract.TransactOpts, withdrawAddress)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address withdrawAddress, uint256 amount) returns()
func (_Contract *ContractTransactor) WithdrawTo(opts *bind.TransactOpts, withdrawAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "withdrawTo", withdrawAddress, amount)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address withdrawAddress, uint256 amount) returns()
func (_Contract *ContractSession) WithdrawTo(withdrawAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.WithdrawTo(&_Contract.TransactOpts, withdrawAddress, amount)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x205c2878.
//
// Solidity: function withdrawTo(address withdrawAddress, uint256 amount) returns()
func (_Contract *ContractTransactorSession) WithdrawTo(withdrawAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.WithdrawTo(&_Contract.TransactOpts, withdrawAddress, amount)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x9e281a98.
//
// Solidity: function withdrawToken(address to, uint256 amount) returns()
func (_Contract *ContractTransactor) WithdrawToken(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "withdrawToken", to, amount)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x9e281a98.
//
// Solidity: function withdrawToken(address to, uint256 amount) returns()
func (_Contract *ContractSession) WithdrawToken(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.WithdrawToken(&_Contract.TransactOpts, to, amount)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x9e281a98.
//
// Solidity: function withdrawToken(address to, uint256 amount) returns()
func (_Contract *ContractTransactorSession) WithdrawToken(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.WithdrawToken(&_Contract.TransactOpts, to, amount)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Contract *ContractTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Contract *ContractSession) Receive() (*types.Transaction, error) {
	return _Contract.Contract.Receive(&_Contract.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Contract *ContractTransactorSession) Receive() (*types.Transaction, error) {
	return _Contract.Contract.Receive(&_Contract.TransactOpts)
}

// ContractConfigUpdatedIterator is returned from FilterConfigUpdated and is used to iterate over the raw logs and unpacked data for ConfigUpdated events raised by the Contract contract.
type ContractConfigUpdatedIterator struct {
	Event *ContractConfigUpdated // Event containing the contract specifics and raw log

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
func (it *ContractConfigUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractConfigUpdated)
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
		it.Event = new(ContractConfigUpdated)
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
func (it *ContractConfigUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractConfigUpdated represents a ConfigUpdated event raised by the Contract contract.
type ContractConfigUpdated struct {
	TokenPaymasterConfig TokenPaymasterTokenPaymasterConfig
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdated is a free log retrieval operation binding the contract event 0xcd938817f1c47094d43be3d07e8c67e11766db2e11a2b4376e7ee937b15793a2.
//
// Solidity: event ConfigUpdated((uint256,uint128,uint48,uint48) tokenPaymasterConfig)
func (_Contract *ContractFilterer) FilterConfigUpdated(opts *bind.FilterOpts) (*ContractConfigUpdatedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ConfigUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractConfigUpdatedIterator{contract: _Contract.contract, event: "ConfigUpdated", logs: logs, sub: sub}, nil
}

// WatchConfigUpdated is a free log subscription operation binding the contract event 0xcd938817f1c47094d43be3d07e8c67e11766db2e11a2b4376e7ee937b15793a2.
//
// Solidity: event ConfigUpdated((uint256,uint128,uint48,uint48) tokenPaymasterConfig)
func (_Contract *ContractFilterer) WatchConfigUpdated(opts *bind.WatchOpts, sink chan<- *ContractConfigUpdated) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ConfigUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractConfigUpdated)
				if err := _Contract.contract.UnpackLog(event, "ConfigUpdated", log); err != nil {
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

// ParseConfigUpdated is a log parse operation binding the contract event 0xcd938817f1c47094d43be3d07e8c67e11766db2e11a2b4376e7ee937b15793a2.
//
// Solidity: event ConfigUpdated((uint256,uint128,uint48,uint48) tokenPaymasterConfig)
func (_Contract *ContractFilterer) ParseConfigUpdated(log types.Log) (*ContractConfigUpdated, error) {
	event := new(ContractConfigUpdated)
	if err := _Contract.contract.UnpackLog(event, "ConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Contract contract.
type ContractOwnershipTransferredIterator struct {
	Event *ContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractOwnershipTransferred)
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
		it.Event = new(ContractOwnershipTransferred)
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
func (it *ContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractOwnershipTransferred represents a OwnershipTransferred event raised by the Contract contract.
type ContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Contract *ContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractOwnershipTransferredIterator{contract: _Contract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Contract *ContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractOwnershipTransferred)
				if err := _Contract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Contract *ContractFilterer) ParseOwnershipTransferred(log types.Log) (*ContractOwnershipTransferred, error) {
	event := new(ContractOwnershipTransferred)
	if err := _Contract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractReceivedIterator is returned from FilterReceived and is used to iterate over the raw logs and unpacked data for Received events raised by the Contract contract.
type ContractReceivedIterator struct {
	Event *ContractReceived // Event containing the contract specifics and raw log

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
func (it *ContractReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractReceived)
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
		it.Event = new(ContractReceived)
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
func (it *ContractReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractReceived represents a Received event raised by the Contract contract.
type ContractReceived struct {
	Sender common.Address
	Value  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterReceived is a free log retrieval operation binding the contract event 0x88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f88525874.
//
// Solidity: event Received(address indexed sender, uint256 value)
func (_Contract *ContractFilterer) FilterReceived(opts *bind.FilterOpts, sender []common.Address) (*ContractReceivedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "Received", senderRule)
	if err != nil {
		return nil, err
	}
	return &ContractReceivedIterator{contract: _Contract.contract, event: "Received", logs: logs, sub: sub}, nil
}

// WatchReceived is a free log subscription operation binding the contract event 0x88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f88525874.
//
// Solidity: event Received(address indexed sender, uint256 value)
func (_Contract *ContractFilterer) WatchReceived(opts *bind.WatchOpts, sink chan<- *ContractReceived, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "Received", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractReceived)
				if err := _Contract.contract.UnpackLog(event, "Received", log); err != nil {
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

// ParseReceived is a log parse operation binding the contract event 0x88a5966d370b9919b20f3e2c13ff65706f196a4e32cc2c12bf57088f88525874.
//
// Solidity: event Received(address indexed sender, uint256 value)
func (_Contract *ContractFilterer) ParseReceived(log types.Log) (*ContractReceived, error) {
	event := new(ContractReceived)
	if err := _Contract.contract.UnpackLog(event, "Received", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractTokenPriceUpdatedIterator is returned from FilterTokenPriceUpdated and is used to iterate over the raw logs and unpacked data for TokenPriceUpdated events raised by the Contract contract.
type ContractTokenPriceUpdatedIterator struct {
	Event *ContractTokenPriceUpdated // Event containing the contract specifics and raw log

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
func (it *ContractTokenPriceUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractTokenPriceUpdated)
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
		it.Event = new(ContractTokenPriceUpdated)
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
func (it *ContractTokenPriceUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractTokenPriceUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractTokenPriceUpdated represents a TokenPriceUpdated event raised by the Contract contract.
type ContractTokenPriceUpdated struct {
	CurrentPrice         *big.Int
	PreviousPrice        *big.Int
	CachedPriceTimestamp *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterTokenPriceUpdated is a free log retrieval operation binding the contract event 0x00d4fe314618b73a96886b87817a53a5ed51433b0234c85a5e9dafe2cb7b8842.
//
// Solidity: event TokenPriceUpdated(uint256 currentPrice, uint256 previousPrice, uint256 cachedPriceTimestamp)
func (_Contract *ContractFilterer) FilterTokenPriceUpdated(opts *bind.FilterOpts) (*ContractTokenPriceUpdatedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "TokenPriceUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractTokenPriceUpdatedIterator{contract: _Contract.contract, event: "TokenPriceUpdated", logs: logs, sub: sub}, nil
}

// WatchTokenPriceUpdated is a free log subscription operation binding the contract event 0x00d4fe314618b73a96886b87817a53a5ed51433b0234c85a5e9dafe2cb7b8842.
//
// Solidity: event TokenPriceUpdated(uint256 currentPrice, uint256 previousPrice, uint256 cachedPriceTimestamp)
func (_Contract *ContractFilterer) WatchTokenPriceUpdated(opts *bind.WatchOpts, sink chan<- *ContractTokenPriceUpdated) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "TokenPriceUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractTokenPriceUpdated)
				if err := _Contract.contract.UnpackLog(event, "TokenPriceUpdated", log); err != nil {
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

// ParseTokenPriceUpdated is a log parse operation binding the contract event 0x00d4fe314618b73a96886b87817a53a5ed51433b0234c85a5e9dafe2cb7b8842.
//
// Solidity: event TokenPriceUpdated(uint256 currentPrice, uint256 previousPrice, uint256 cachedPriceTimestamp)
func (_Contract *ContractFilterer) ParseTokenPriceUpdated(log types.Log) (*ContractTokenPriceUpdated, error) {
	event := new(ContractTokenPriceUpdated)
	if err := _Contract.contract.UnpackLog(event, "TokenPriceUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractUniswapRevertedIterator is returned from FilterUniswapReverted and is used to iterate over the raw logs and unpacked data for UniswapReverted events raised by the Contract contract.
type ContractUniswapRevertedIterator struct {
	Event *ContractUniswapReverted // Event containing the contract specifics and raw log

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
func (it *ContractUniswapRevertedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractUniswapReverted)
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
		it.Event = new(ContractUniswapReverted)
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
func (it *ContractUniswapRevertedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractUniswapRevertedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractUniswapReverted represents a UniswapReverted event raised by the Contract contract.
type ContractUniswapReverted struct {
	TokenIn      common.Address
	TokenOut     common.Address
	AmountIn     *big.Int
	AmountOutMin *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterUniswapReverted is a free log retrieval operation binding the contract event 0xf7edd4c6ec425decf715a8b8eaa3b65d3d86e31ad0ff750aa60fa834190f515f.
//
// Solidity: event UniswapReverted(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin)
func (_Contract *ContractFilterer) FilterUniswapReverted(opts *bind.FilterOpts) (*ContractUniswapRevertedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "UniswapReverted")
	if err != nil {
		return nil, err
	}
	return &ContractUniswapRevertedIterator{contract: _Contract.contract, event: "UniswapReverted", logs: logs, sub: sub}, nil
}

// WatchUniswapReverted is a free log subscription operation binding the contract event 0xf7edd4c6ec425decf715a8b8eaa3b65d3d86e31ad0ff750aa60fa834190f515f.
//
// Solidity: event UniswapReverted(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin)
func (_Contract *ContractFilterer) WatchUniswapReverted(opts *bind.WatchOpts, sink chan<- *ContractUniswapReverted) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "UniswapReverted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractUniswapReverted)
				if err := _Contract.contract.UnpackLog(event, "UniswapReverted", log); err != nil {
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

// ParseUniswapReverted is a log parse operation binding the contract event 0xf7edd4c6ec425decf715a8b8eaa3b65d3d86e31ad0ff750aa60fa834190f515f.
//
// Solidity: event UniswapReverted(address tokenIn, address tokenOut, uint256 amountIn, uint256 amountOutMin)
func (_Contract *ContractFilterer) ParseUniswapReverted(log types.Log) (*ContractUniswapReverted, error) {
	event := new(ContractUniswapReverted)
	if err := _Contract.contract.UnpackLog(event, "UniswapReverted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractUserOperationSponsoredIterator is returned from FilterUserOperationSponsored and is used to iterate over the raw logs and unpacked data for UserOperationSponsored events raised by the Contract contract.
type ContractUserOperationSponsoredIterator struct {
	Event *ContractUserOperationSponsored // Event containing the contract specifics and raw log

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
func (it *ContractUserOperationSponsoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractUserOperationSponsored)
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
		it.Event = new(ContractUserOperationSponsored)
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
func (it *ContractUserOperationSponsoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractUserOperationSponsoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractUserOperationSponsored represents a UserOperationSponsored event raised by the Contract contract.
type ContractUserOperationSponsored struct {
	User                       common.Address
	ActualTokenCharge          *big.Int
	ActualGasCost              *big.Int
	ActualTokenPriceWithMarkup *big.Int
	Raw                        types.Log // Blockchain specific contextual infos
}

// FilterUserOperationSponsored is a free log retrieval operation binding the contract event 0x46caa0511cf037f06f57a0bf273a2ff04229f5b12fb04675234a6cbe2e7f1a89.
//
// Solidity: event UserOperationSponsored(address indexed user, uint256 actualTokenCharge, uint256 actualGasCost, uint256 actualTokenPriceWithMarkup)
func (_Contract *ContractFilterer) FilterUserOperationSponsored(opts *bind.FilterOpts, user []common.Address) (*ContractUserOperationSponsoredIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "UserOperationSponsored", userRule)
	if err != nil {
		return nil, err
	}
	return &ContractUserOperationSponsoredIterator{contract: _Contract.contract, event: "UserOperationSponsored", logs: logs, sub: sub}, nil
}

// WatchUserOperationSponsored is a free log subscription operation binding the contract event 0x46caa0511cf037f06f57a0bf273a2ff04229f5b12fb04675234a6cbe2e7f1a89.
//
// Solidity: event UserOperationSponsored(address indexed user, uint256 actualTokenCharge, uint256 actualGasCost, uint256 actualTokenPriceWithMarkup)
func (_Contract *ContractFilterer) WatchUserOperationSponsored(opts *bind.WatchOpts, sink chan<- *ContractUserOperationSponsored, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "UserOperationSponsored", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractUserOperationSponsored)
				if err := _Contract.contract.UnpackLog(event, "UserOperationSponsored", log); err != nil {
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

// ParseUserOperationSponsored is a log parse operation binding the contract event 0x46caa0511cf037f06f57a0bf273a2ff04229f5b12fb04675234a6cbe2e7f1a89.
//
// Solidity: event UserOperationSponsored(address indexed user, uint256 actualTokenCharge, uint256 actualGasCost, uint256 actualTokenPriceWithMarkup)
func (_Contract *ContractFilterer) ParseUserOperationSponsored(log types.Log) (*ContractUserOperationSponsored, error) {
	event := new(ContractUserOperationSponsored)
	if err := _Contract.contract.UnpackLog(event, "UserOperationSponsored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
