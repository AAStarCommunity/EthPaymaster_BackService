package userop

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/paymaster_abi"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/xerrors"
	"math/big"
	"reflect"
	"strconv"
	"sync"
)

var (
	validate              = validator.New()
	onlyOnce              = sync.Once{}
	userOPV1HashArguments abi.Arguments
	userOpV1PackArg       abi.Arguments
)

func init() {
	userOPV1HashArguments = abi.Arguments{
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
	userOpV1PackArg := abi.Arguments{
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
}

func NewUserOp(userOp *map[string]any, entryPointVersion types.EntrypointVersion) (*BaseUserOp, error) {
	var result BaseUserOp
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
	// Validate struct
	err = validate.Struct(result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type BaseUserOp interface {
	GetEntrypointVersion() types.EntrypointVersion

	GetUserOpHash(strategy *model.Strategy) ([]byte, string, error)
	GetSender() *common.Address
	Pack() (string, []byte, error)
	ValidateUserOp() error
	GetCallData() []byte
}
type BaseUserOperation struct {
	Sender               *common.Address `json:"sender"   mapstructure:"sender"  binding:"required,hexParam"`
	InitCode             []byte          `json:"initCode"  mapstructure:"init_code" `
	CallData             []byte          `json:"callData"  mapstructure:"call_data"  binding:"required"`
	PreVerificationGas   *big.Int        `json:"preVerificationGas"  mapstructure:"pre_verification_gas"  binding:"required"`
	MaxFeePerGas         *big.Int        `json:"maxFeePerGas"  mapstructure:"max_fee_per_gas"  binding:"required"`
	PaymasterAndData     []byte          `json:"paymasterAndData"  mapstructure:"paymaster_and_data"`
	Signature            []byte          `json:"signature"  mapstructure:"signature"  binding:"required"`
	Nonce                *big.Int        `json:"nonce"  mapstructure:"nonce"  binding:"required"`
	MaxPriorityFeePerGas *big.Int        `json:"maxPriorityFeePerGas"  mapstructure:"max_priority_fee_per_gas"  binding:"required"`
}

// UserOperation  entrypoint v0.0.6
// verificationGasLimit validateUserOp ,validatePaymasterUserOp limit
// callGasLimit calldata Execute gas limit
// preVerificationGas
type UserOperation struct {
	BaseUserOperation
	CallGasLimit         *big.Int `json:"callGasLimit"  mapstructure:"call_gas_limit"  binding:"required"`
	VerificationGasLimit *big.Int `json:"verificationGasLimit"  mapstructure:"verification_gas_limit"  binding:"required"`
}

func (userOp *UserOperation) GetEntrypointVersion() types.EntrypointVersion {
	return types.EntrypointV06
}
func (userOp *UserOperation) GetSender() *common.Address {
	return userOp.Sender
}
func (userOp *UserOperation) ValidateUserOp() error {
	return nil
}
func (userOp *UserOperation) GetCallData() []byte {
	return userOp.CallData
}

func (userOp *UserOperation) GetUserOpHash(strategy *model.Strategy) ([]byte, string, error) {
	packUserOpStr, _, err := userOp.Pack()
	if err != nil {
		return nil, "", err
	}

	chainId := conf.GetChainId(strategy.GetNewWork())
	if chainId == nil {
		return nil, "", xerrors.Errorf("empty NetWork")
	}
	if err != nil {
		return nil, "", err
	}
	packUserOpStrByteNew, _ := hex.DecodeString(packUserOpStr)

	bytesRes, err := userOPV1HashArguments.Pack(packUserOpStrByteNew, chainId, strategy.GetPaymasterAddress(), userOp.Nonce, strategy.ExecuteRestriction.EffectiveStartTime, strategy.ExecuteRestriction.EffectiveEndTime)
	if err != nil {
		return nil, "", err
	}

	encodeHash := crypto.Keccak256(bytesRes)
	return encodeHash, hex.EncodeToString(bytesRes), nil
}

func (userOp *UserOperation) Pack() (string, []byte, error) {
	//TODO disgusting logic
	paymasterDataMock := "d93349Ee959d295B115Ee223aF10EF432A8E8523000000000000000000000000000000000000000000000000000000001710044496000000000000000000000000000000000000000000000000000000174158049605bea0bfb8539016420e76749fda407b74d3d35c539927a45000156335643827672fa359ee968d72db12d4b4768e8323cd47443505ab138a525c1f61c6abdac501"
	encoded, err := userOpV1PackArg.Pack(userOp.Sender, userOp.Nonce, userOp.InitCode, userOp.CallData, userOp.CallGasLimit, userOp.VerificationGasLimit, userOp.PreVerificationGas, userOp.MaxFeePerGas, userOp.MaxPriorityFeePerGas, paymasterDataMock, userOp.Sender)
	//fmt.Printf("paymasterDataTmpLen: %x\n", len(paymasterDataTmp))
	//fmt.Printf("paymasterDataKLen : %x\n", len(userOp.PaymasterAndData))

	if err != nil {
		return "", nil, err
	}
	//https://github.com/jayden-sudo/SoulWalletCore/blob/dc76bdb9a156d4f99ef41109c59ab99106c193ac/contracts/utils/CalldataPack.sol#L51-L65
	hexString := hex.EncodeToString(encoded)
	//1. 从 63*10+ 1 ～64*10获取
	hexString = hexString[64:]
	//hexLen := len(hexString)
	subIndex := GetIndex(hexString)
	hexString = hexString[:subIndex]
	//fmt.Printf("subIndex: %d\n", subIndex)
	return hexString, encoded, nil
}

/**
Userop V2
**/

// UserOperationV2  entrypoint v0.0.7
type UserOperationV2 struct {
	BaseUserOperation
	AccountGasLimit string `json:"account_gas_limit" binding:"required"`
}

func (u *UserOperationV2) GetEntrypointVersion() types.EntrypointVersion {
	return types.EntryPointV07
}
func (userOp *UserOperationV2) ValidateUserOp() error {
	return nil

}
func (u *UserOperationV2) GetSender() *common.Address {
	return u.Sender
}
func (userOp *UserOperationV2) GetCallData() []byte {
	return userOp.CallData
}

func (userOp *UserOperationV2) Pack() (string, []byte, error) {
	return "", nil, nil
}

func (userOp *UserOperationV2) GetUserOpHash(strategy *model.Strategy) ([]byte, string, error) {

	return nil, "", nil
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
