package user_op

import (
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_contract/paymaster_abi"
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/xerrors"
	"math/big"
	"reflect"
	"strconv"
	"sync"
)

var (
	DummyGasFees              = utils.PackIntTo32Bytes(global_const.DummyMaxPriorityFeePerGas, global_const.DummyMaxFeePerGas)
	DummyAccountGasLimits     = utils.PackIntTo32Bytes(global_const.DummyVerificationGasLimit, global_const.DummyCallGasLimit)
	MinPreVerificationGas     *big.Int
	validate                  = validator.New()
	onlyOnce                  = sync.Once{}
	userOPV06GetHashArguments abi.Arguments
	UserOpV07GetHashArguments abi.Arguments
	userOpV06PackArg          abi.Arguments
	UserOpV07PackArg          abi.Arguments
)

func init() {
	MinPreVerificationGas = big.NewInt(21000)

	userOPV06GetHashArguments = abi.Arguments{
		{
			Type: paymaster_abi.BytesType,
		},
		{
			Type: paymaster_abi.Uint256Type,
		},
		{
			Type: paymaster_abi.AddressType,
		},
		{
			Type: paymaster_abi.Uint256Type,
		},
		{
			Type: paymaster_abi.Uint48Type,
		},
		{
			Type: paymaster_abi.Uint48Type,
		},
	}
	UserOpV07GetHashArguments = abi.Arguments{
		{
			Type: paymaster_abi.AddressType, //Sender
		},
		{
			Type: paymaster_abi.Uint256Type, //userOp.nonce
		},
		{
			Type: paymaster_abi.Bytes32Type, //keccak256(userOp.initCode),
		},
		{
			Type: paymaster_abi.Bytes32Type, //keccak256(userOp.callData),
		},
		{
			Type: paymaster_abi.Bytes32Type, //userOp.accountGasLimits,
		},
		{
			Type: paymaster_abi.Uint256Type, //uint256(bytes32(userOp.paymasterAndData[PAYMASTER_VALIDATION_GAS_OFFSET : PAYMASTER_DATA_OFFSET])),
		},
		{
			Type: paymaster_abi.Uint256Type, // userOp.preVerificationGas,
		},
		{
			Type: paymaster_abi.Uint256Type, //userOp.gasFees,
		},
		{
			Type: paymaster_abi.Uint256Type, //block.chainid,
		},
		{
			Type: paymaster_abi.AddressType, //address(this),
		},
		{
			Type: paymaster_abi.Uint48Type, // validUntil,
		},
		{
			Type: paymaster_abi.Uint48Type, // validAfter,
		},
	}
	userOpV06PackArg = abi.Arguments{
		{
			Name: "Sender",
			Type: paymaster_abi.AddressType,
		},
		{
			Name: "Nonce",
			Type: paymaster_abi.Uint256Type,
		},
		{
			Name: "InitCode",
			Type: paymaster_abi.BytesType,
		},
		{
			Name: "CallData",
			Type: paymaster_abi.BytesType,
		},
		{
			Name: "CallGasLimit",
			Type: paymaster_abi.Uint256Type,
		},
		{
			Name: "VerificationGasLimit",
			Type: paymaster_abi.Uint256Type,
		},
		{
			Name: "PreVerificationGas",
			Type: paymaster_abi.Uint256Type,
		},
		{
			Name: "MaxFeePerGas",
			Type: paymaster_abi.Uint256Type,
		},
		{
			Name: "MaxPriorityFeePerGas",
			Type: paymaster_abi.Uint256Type,
		},
		{
			Name: "PaymasterAndData",
			Type: paymaster_abi.BytesType,
		},
		{
			Name: "Signature",
			Type: paymaster_abi.BytesType,
		},
	}
	UserOpV07PackArg = abi.Arguments{
		{
			Name: "Sender",
			Type: paymaster_abi.AddressType,
		},
		{
			Name: "Nonce",
			Type: paymaster_abi.Uint256Type,
		},
		{
			Name: "InitCode",
			Type: paymaster_abi.BytesType,
		},
		{
			Name: "CallData",
			Type: paymaster_abi.BytesType,
		},
		{
			Name: "AccountGasLimits",
			Type: paymaster_abi.Bytes32Type,
		},
		{
			Name: "PreVerificationGas",
			Type: paymaster_abi.Uint256Type,
		},
		{
			Name: "GasFees",
			Type: paymaster_abi.Bytes32Type,
		},
		{
			Name: "PaymasterAndData",
			Type: paymaster_abi.BytesType,
		},
		{
			Name: "Signature",
			Type: paymaster_abi.BytesType,
		},
	}
}

type BaseUserOperation struct {
	Sender             *common.Address `json:"sender"   mapstructure:"sender"  binding:"required,hexParam"`
	Nonce              *big.Int        `json:"nonce"  mapstructure:"nonce"  binding:"required"`
	InitCode           []byte          `json:"initCode"  mapstructure:"initCode" `
	CallData           []byte          `json:"callData"  mapstructure:"callData"  binding:"required"`
	PreVerificationGas *big.Int        `json:"preVerificationGas"  mapstructure:"pre_verification_gas"  binding:"required"`
	PaymasterAndData   []byte          `json:"paymasterAndData"  mapstructure:"paymaster_and_data"`
	Signature          []byte          `json:"signature"  mapstructure:"signature"  binding:"required"`
}

// UserOperationV06  entrypoint v0.0.6
// verificationGasLimit validateUserOp ,validatePaymasterUserOp limit
// callGasLimit calldata Execute gas limit
// preVerificationGas
type UserOperationV06 struct {
	Sender             *common.Address `json:"sender"   mapstructure:"sender"  binding:"required,hexParam"`
	Nonce              *big.Int        `json:"nonce"  mapstructure:"nonce"  binding:"required"`
	InitCode           []byte          `json:"initCode"  mapstructure:"initCode" `
	CallData           []byte          `json:"callData"  mapstructure:"callData"  binding:"required"`
	PreVerificationGas *big.Int        `json:"preVerificationGas"  mapstructure:"pre_verification_gas"  binding:"required"`
	PaymasterAndData   []byte          `json:"paymasterAndData"  mapstructure:"paymaster_and_data"`
	Signature          []byte          `json:"signature"  mapstructure:"signature"  binding:"required"`
	//Maximum fee per gas (similar to EIP-1559  max_fee_per_gas)
	MaxFeePerGas *big.Int `json:"maxFeePerGas"  mapstructure:"maxFeePerGas"  binding:"required"`
	//Maximum priority fee per gas (similar to EIP-1559 max_priority_fee_per_gas)
	MaxPriorityFeePerGas *big.Int `json:"maxPriorityFeePerGas"  mapstructure:"max_priority_fee_per_gas"  binding:"required"`
	//Gas limit for execution phase
	CallGasLimit *big.Int `json:"callGasLimit"  mapstructure:"call_gas_limit"  binding:"required"`
	//Gas limit for verification phase
	VerificationGasLimit *big.Int `json:"verificationGasLimit"  mapstructure:"verification_gas_limit"  binding:"required"`
}

// UserOperationV07  entrypoint v0.0.7
type UserOperationV07 struct {
	BaseUserOperation
	AccountGasLimit               [32]byte `json:"account_gas_limit" mapstructure:"account_gas_limit" binding:"required"`
	PaymasterVerificationGasLimit *big.Int `json:"paymaster_verification_gas_limit" mapstructure:"paymaster_verification_gas_limit" binding:"required"`
	PaymasterPostOpGasLimit       *big.Int `json:"paymaster_post_op_gas_limit" mapstructure:"paymaster_post_op_gas_limit" binding:"required"`
	GasFees                       [32]byte `json:"gas_fees" mapstructure:"gas_fees" binding:"required"`
}

func NewUserOp(userOp *map[string]any) (*UserOpInput, error) {
	var result UserOpInput

	// Convert map to struct
	decodeConfig := &mapstructure.DecoderConfig{
		DecodeHook: decodeOpTypes,
		Result:     &result,
		ErrorUnset: false,
		MatchName:  exactFieldMatch,
	}
	decoder, err := mapstructure.NewDecoder(decodeConfig)
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(userOp); err != nil {
		return nil, xerrors.Errorf("data [%w] convert failed: [%w]", userOp, err)
	}
	onlyOnce.Do(func() {
		validate.RegisterCustomTypeFunc(validateAddressType, common.Address{})
		validate.RegisterCustomTypeFunc(validateBigIntType, big.Int{})
	})
	// Validate struct
	err = validate.Struct(result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type UserOpInput struct {
	Sender             *common.Address `json:"sender"   mapstructure:"sender"  binding:"required,hexParam"`
	Nonce              *big.Int        `json:"nonce"  mapstructure:"nonce"  binding:"required"`
	InitCode           []byte          `json:"initCode"  mapstructure:"initCode" `
	CallData           []byte          `json:"callData"  mapstructure:"callData"  binding:"required"`
	PreVerificationGas *big.Int        `json:"preVerificationGas"  mapstructure:"preVerificationGas"`
	PaymasterAndData   []byte          `json:"paymasterAndData"  mapstructure:"paymasterAndData"`
	Signature          []byte          `json:"signature"  mapstructure:"signature"  binding:"required"`
	//Maximum fee per gas (similar to EIP-1559  max_fee_per_gas)
	MaxFeePerGas *big.Int `json:"maxFeePerGas"  mapstructure:"maxFeePerGas"`
	//Maximum priority fee per gas (similar to EIP-1559 max_priority_fee_per_gas)
	MaxPriorityFeePerGas *big.Int `json:"maxPriorityFeePerGas"  mapstructure:"maxPriorityFeePerGas" `
	//Gas limit for execution phase
	CallGasLimit *big.Int `json:"callGasLimit"  mapstructure:"callGasLimit"`
	//Gas limit for verification phase
	VerificationGasLimit *big.Int `json:"verificationGasLimit"  mapstructure:"verificationGasLimit"`
	AccountGasLimits     [32]byte `json:"accountGasLimits"  mapstructure:"accountGasLimits"`
	GasFees              [32]byte `json:"gasFees"  mapstructure:"gas_fees"  binding:"required"`

	ComputeGasOnly bool
}

func (userOp *UserOpInput) PackUserOpForMock(version global_const.EntrypointVersion) (string, []byte, error) {
	if version == global_const.EntrypointV07 {
		gasFee := utils.PackIntTo32Bytes(userOp.MaxPriorityFeePerGas, userOp.MaxFeePerGas)
		encoded, err := UserOpV07PackArg.Pack(userOp.Sender, userOp.Nonce, userOp.InitCode, userOp.CallData, DummyAccountGasLimits, userOp.PreVerificationGas, gasFee, global_const.DummyPaymasterDataByte, global_const.DummySignatureByte)
		if err != nil {
			return "", nil, err
		}
		return hex.EncodeToString(encoded), encoded, nil
	} else if version == global_const.EntrypointV06 {
		maxFeePerGas := userOp.MaxFeePerGas
		if maxFeePerGas == nil {
			maxFeePerGas = global_const.DummyMaxFeePerGas
		}
		maxPriorityFeePerGas := userOp.MaxPriorityFeePerGas
		if maxPriorityFeePerGas == nil {
			maxPriorityFeePerGas = global_const.DummyMaxPriorityFeePerGas
		}
		encoded, err := userOpV06PackArg.Pack(userOp.Sender, userOp.Nonce, userOp.InitCode, userOp.CallData, global_const.DummyCallGasLimit, global_const.DummyVerificationGasLimit, global_const.DummyReverificationsBigint, maxFeePerGas, maxPriorityFeePerGas, global_const.DummyPaymasterDataByte, userOp.Signature)

		if err != nil {
			return "", nil, err
		}
		return hex.EncodeToString(encoded), encoded, nil
	} else {
		return "", nil, xerrors.Errorf("should never be here ")
	}
}

func (userOp *UserOpInput) GetFactoryAddress() (*common.Address, error) {
	initCode := userOp.InitCode
	if len(initCode) < 20 {
		return nil, xerrors.Errorf("initCode length is less than 20")
	}
	address := common.BytesToAddress(initCode[:42])
	return &address, nil
}

func (userOp *UserOpInput) ValidateUserOp() error {
	panic("implement me")
}

func GetIndex(hexString string) int64 {
	//1. 从 63*10+ 1 ～64*10获取
	indexPre := hexString[576:640]
	indePreInt, _ := strconv.ParseInt(indexPre, 16, 64)
	result := indePreInt * 2
	return result
}
func validateAddressType(field reflect.Value) interface{} {
	value, ok := field.Interface().(common.Address)
	if !ok || value == common.HexToAddress("0x") {
		return nil
	}
	return field
}

func validateBigIntType(field reflect.Value) interface{} {
	value, ok := field.Interface().(big.Int)
	if !ok || value.Cmp(big.NewInt(0)) == -1 {
		return nil
	}
	return field
}

func exactFieldMatch(mapKey, fieldName string) bool {
	return mapKey == fieldName
}

func decodeOpTypes(
	f reflect.Kind,
	t reflect.Kind,
	data interface{}) (interface{}, error) {
	// String to common.Address conversion
	if f == reflect.String && t == reflect.Array {
		return common.HexToAddress(data.(string)), nil
	}

	// String to big.Int conversion
	if f == reflect.String && t == reflect.Struct {
		n := new(big.Int)
		n, ok := n.SetString(data.(string), 0)
		if !ok {
			return nil, xerrors.Errorf("bigInt conversion failed")
		}
		return n, nil
	}

	// Float64 to big.Int conversion
	if f == reflect.Float64 && t == reflect.Struct {
		n, ok := data.(float64)
		if !ok {
			return nil, xerrors.Errorf("bigInt conversion failed")
		}
		return big.NewInt(int64(n)), nil
	}

	// String to []byte conversion
	if f == reflect.String && t == reflect.Slice {
		byteStr := data.(string)
		if len(byteStr) < 2 || byteStr[:2] != "0x" {
			return nil, xerrors.Errorf("not byte string")
		}

		b, err := hex.DecodeString(byteStr[2:])
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	return data, nil
}
