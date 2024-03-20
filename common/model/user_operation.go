package model

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/xerrors"
	"math/big"
	"reflect"
	"sync"
)

var (
	validate = validator.New()
	onlyOnce = sync.Once{}
)

// UserOperation  entrypoint v0.0.6
// verificationGasLimit validateUserOp ,validatePaymasterUserOp limit
// callGasLimit calldata Execute gas limit
// preVerificationGas
type UserOperation struct {
	Sender               common.Address `json:"sender"   mapstructure:"sender"  binding:"required,hexParam"`
	Nonce                *big.Int       `json:"nonce"  mapstructure:"nonce"  binding:"required"`
	InitCode             []byte         `json:"initCode"  mapstructure:"init_code" `
	CallData             []byte         `json:"callData"  mapstructure:"call_data"  binding:"required"`
	CallGasLimit         *big.Int       `json:"callGasLimit"  mapstructure:"call_gas_limit"  binding:"required"`
	VerificationGasLimit *big.Int       `json:"verificationGasLimit"  mapstructure:"verification_gas_limit"  binding:"required"`
	PreVerificationGas   *big.Int       `json:"preVerificationGas"  mapstructure:"pre_verification_gas"  binding:"required"`
	MaxFeePerGas         *big.Int       `json:"maxFeePerGas"  mapstructure:"max_fee_per_gas"  binding:"required"`
	MaxPriorityFeePerGas *big.Int       `json:"maxPriorityFeePerGas"  mapstructure:"max_priority_fee_per_gas"  binding:"required"`
	PaymasterAndData     []byte         `json:"paymasterAndData"  mapstructure:"paymaster_and_data"`
	Signature            []byte         `json:"signature"  mapstructure:"signature"  binding:"required"`
}
type UserOperationSimple struct {
	Sender []byte `json:"sender"   mapstructure:"sender"  binding:"required,hexParam"`
}

// PackUserOperation  entrypoint v0.0.67
type PackUserOperation struct {
	Sender               string `json:"sender" binding:"required,hexParam"`
	Nonce                string `json:"nonce" binding:"required"`
	InitCode             string `json:"init_code"`
	CallData             string `json:"call_data" binding:"required"`
	AccountGasLimit      string `json:"account_gas_limit" binding:"required"`
	PreVerificationGas   string `json:"pre_verification_gas" binding:"required"`
	MaxFeePerGas         string `json:"max_fee_per_gas" binding:"required"`
	MaxPriorityFeePerGas string `json:"max_priority_fee_per_gas" binding:"required"`
	PaymasterAndData     string `json:"paymaster_and_data"`
	Signature            string `json:"signature" binding:"required"`
}

func NewUserOp(userOp *map[string]any) (*UserOperation, error) {
	var result UserOperation
	// Convert map to struct
	decodeConfig := &mapstructure.DecoderConfig{
		DecodeHook: decodeOpTypes,
		Result:     &result,
		ErrorUnset: true,
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
	err = validate.Struct(result)
	if err != nil {
		return nil, err
	}

	return &result, nil
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
