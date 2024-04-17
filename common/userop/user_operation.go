package userop

import (
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_common/paymaster_abi"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
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
			Type: paymaster_abi.Uint256Type,
		},
		{
			Name: "PreVerificationGas",
			Type: paymaster_abi.Uint256Type,
		},
		{
			Name: "GasFees",
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
	if entryPointVersion == types.EntryPointV07 {
		result = &UserOperationV07{}
	} else {
		result = &UserOperationV06{}
	}
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
	PackUserOpForMock() (string, []byte, error)
	ValidateUserOp() error
	GetInitCode() []byte
	GetCallData() []byte
	GetNonce() *big.Int
	GetFactoryAddress() *common.Address
	SetSignature(signature []byte)
	SetPreVerificationGas(preVerificationGas *big.Int)
}
type BaseUserOperation struct {
	Sender             *common.Address `json:"sender"   mapstructure:"sender"  binding:"required,hexParam"`
	Nonce              *big.Int        `json:"nonce"  mapstructure:"nonce"  binding:"required"`
	InitCode           []byte          `json:"initCode"  mapstructure:"init_code" `
	CallData           []byte          `json:"callData"  mapstructure:"call_data"  binding:"required"`
	PreVerificationGas *big.Int        `json:"preVerificationGas"  mapstructure:"pre_verification_gas"  binding:"required"`
	PaymasterAndData   []byte          `json:"paymasterAndData"  mapstructure:"paymaster_and_data"`
	Signature          []byte          `json:"signature"  mapstructure:"signature"  binding:"required"`
}

// UserOperationV06  entrypoint v0.0.6
// verificationGasLimit validateUserOp ,validatePaymasterUserOp limit
// callGasLimit calldata Execute gas limit
// preVerificationGas
type UserOperationV06 struct {
	BaseUserOperation
	//Maximum fee per gas (similar to EIP-1559  max_fee_per_gas)
	MaxFeePerGas *big.Int `json:"maxFeePerGas"  mapstructure:"max_fee_per_gas"  binding:"required"`
	//Maximum priority fee per gas (similar to EIP-1559 max_priority_fee_per_gas)
	MaxPriorityFeePerGas *big.Int `json:"maxPriorityFeePerGas"  mapstructure:"max_priority_fee_per_gas"  binding:"required"`
	//Gas limit for execution phase
	CallGasLimit *big.Int `json:"callGasLimit"  mapstructure:"call_gas_limit"  binding:"required"`
	//Gas limit for verification phase
	VerificationGasLimit *big.Int `json:"verificationGasLimit"  mapstructure:"verification_gas_limit"  binding:"required"`
}

func (userOp *UserOperationV06) SetSignature(signature []byte) {
	userOp.Signature = signature
}
func (userOp *UserOperationV06) SetPreVerificationGas(preVerificationGas *big.Int) {
	userOp.PreVerificationGas = preVerificationGas
}
func (userOp *UserOperationV06) GetEntrypointVersion() types.EntrypointVersion {
	return types.EntrypointV06
}
func (userOp *UserOperationV06) GetSender() *common.Address {
	return userOp.Sender
}
func (userOp *UserOperationV06) ValidateUserOp() error {
	if userOp.PreVerificationGas.Cmp(MinPreVerificationGas) == -1 {
		return xerrors.Errorf("preVerificationGas is less than 21000")
	}
	return nil
}
func (userOp *UserOperationV06) GetCallData() []byte {
	return userOp.CallData
}
func (userop *UserOperationV06) GetInitCode() []byte {
	return userop.InitCode
}
func (userOp *UserOperationV06) GetNonce() *big.Int {
	return userOp.Nonce
}
func (userOp *UserOperationV06) GetFactoryAddress() *common.Address {
	//TODO
	return nil
}

// PackUserOpForMock return keccak256(abi.encode(
//
//	    pack(userOp),
//	    block.chainid,
//	    address(this),
//	    senderNonce[userOp.getSender()],
//	    validUntil,
//	    validAfter
//	));
func packUserOpV6ForUserOpHash(userOp *UserOperationV06) (string, []byte, error) {
	//TODO disgusting logic
	encoded, err := userOpV06PackArg.Pack(userOp.Sender, userOp.Nonce, userOp.InitCode, userOp.CallData, userOp.CallGasLimit, userOp.VerificationGasLimit, userOp.PreVerificationGas, userOp.MaxFeePerGas, userOp.MaxPriorityFeePerGas, types.DUMMY_PAYMASTER_DATA, userOp.Sender)
	if err != nil {
		return "", nil, err
	}
	//https://github.com/jayden-sudo/SoulWalletCore/blob/dc76bdb9a156d4f99ef41109c59ab99106c193ac/contracts/utils/CalldataPack.sol#L51-L65
	hexString := hex.EncodeToString(encoded)
	//1. get From  63*10+ 1 ～64*10
	hexString = hexString[64:]
	//hexLen := len(hexString)
	subIndex := GetIndex(hexString)
	hexString = hexString[:subIndex]
	//fmt.Printf("subIndex: %d\n", subIndex)
	return hexString, encoded, nil
}
func (userOp *UserOperationV06) GetUserOpHash(strategy *model.Strategy) ([]byte, string, error) {
	packUserOpStr, _, err := packUserOpV6ForUserOpHash(userOp)
	if err != nil {
		return nil, "", err
	}

	packUserOpStrByteNew, err := hex.DecodeString(packUserOpStr)
	if err != nil {
		return nil, "", err
	}

	bytesRes, err := userOPV06GetHashArguments.Pack(packUserOpStrByteNew, conf.GetChainId(strategy.GetNewWork()), strategy.GetPaymasterAddress(), userOp.Nonce, strategy.ExecuteRestriction.EffectiveStartTime, strategy.ExecuteRestriction.EffectiveEndTime)
	if err != nil {
		return nil, "", err
	}

	userOpHash := crypto.Keccak256(bytesRes)
	afterProcessUserOphash := utils.ToEthSignedMessageHash(userOpHash)
	return afterProcessUserOphash, hex.EncodeToString(bytesRes), nil

}

func (userOp *UserOperationV06) PackUserOpForMock() (string, []byte, error) {
	//TODO  UserMock
	encoded, err := userOpV06PackArg.Pack(userOp.Sender, userOp.Nonce, userOp.InitCode, userOp.CallData, userOp.CallGasLimit, userOp.VerificationGasLimit, userOp.PreVerificationGas, userOp.MaxFeePerGas, userOp.MaxPriorityFeePerGas, types.DUMMY_PAYMASTER_DATA, userOp.Sender)
	if err != nil {
		return "", nil, err
	}
	return hex.EncodeToString(encoded), encoded, nil
}

// return
func (userOp *UserOperationV06) EstimateGasLimit(strategy *model.Strategy) (uint64, uint64, error) {
	return 0, 0, nil
}
func getVerificationGasLimit(strategy *model.Strategy, userOp *UserOperationV06) (uint64, error) {
	return 0, nil
}

/**
Userop V2
**/

// UserOperationV07  entrypoint v0.0.7
type UserOperationV07 struct {
	BaseUserOperation
	AccountGasLimit               *big.Int `json:"account_gas_limit" binding:"required"`
	PaymasterVerificationGasLimit *big.Int `json:"paymaster_verification_gas_limit" binding:"required"`
	PaymasterPostOpGasLimit       *big.Int `json:"paymaster_post_op_gas_limit" binding:"required"`
	GasFees                       []byte   `json:"gasFees" binding:"required"`
}

func (userOp *UserOperationV07) GetEntrypointVersion() types.EntrypointVersion {
	return types.EntryPointV07
}
func (userOp *UserOperationV07) ValidateUserOp() error {
	return nil

}

func (userOp *UserOperationV07) SetSignature(signature []byte) {
	userOp.Signature = signature
}
func (userOp *UserOperationV07) SetPreVerificationGas(preVerificationGas *big.Int) {
	userOp.PreVerificationGas = preVerificationGas
}
func (userOp *UserOperationV07) GetFactoryAddress() *common.Address {
	//TODO
	return nil
}
func (userOp *UserOperationV07) GetInitCode() []byte {
	return userOp.InitCode
}
func (userOp *UserOperationV07) GetSender() *common.Address {
	return userOp.Sender
}
func (userOp *UserOperationV07) GetCallData() []byte {
	return userOp.CallData
}
func (userOp *UserOperationV07) GetNonce() *big.Int {
	return userOp.Nonce
}
func (userOp *UserOperationV07) PackUserOpForMock() (string, []byte, error) {
	encoded, err := UserOpV07PackArg.Pack(userOp.Sender, userOp.Nonce, userOp.InitCode, userOp.CallData, userOp.AccountGasLimit, userOp.PreVerificationGas, userOp.PreVerificationGas, userOp.GasFees, types.DUMMY_PAYMASTER_DATA, userOp.Signature)
	if err != nil {
		return "", nil, err
	}
	return hex.EncodeToString(encoded), encoded, nil
}

// GetUserOpHash return keccak256(
//
//	abi.encode(
//	    sender,
//	    userOp.nonce,
//	    keccak256(userOp.initCode),
//	    keccak256(userOp.callData),
//	    userOp.accountGasLimits,
//	    uint256(bytes32(userOp.paymasterAndData[PAYMASTER_VALIDATION_GAS_OFFSET : PAYMASTER_DATA_OFFSET])),
//	    userOp.preVerificationGas,
//	    userOp.gasFees,
//	    block.chainid,
//	    address(this)
//	)
func (userOp *UserOperationV07) GetUserOpHash(strategy *model.Strategy) ([]byte, string, error) {
	//TODO
	paymasterGasValue := userOp.PaymasterPostOpGasLimit.Text(20) + userOp.PaymasterVerificationGasLimit.Text(20)
	byteRes, err := UserOpV07GetHashArguments.Pack(userOp.Sender, userOp.Nonce, crypto.Keccak256(userOp.InitCode),
		crypto.Keccak256(userOp.CallData), userOp.AccountGasLimit,
		paymasterGasValue, userOp.PreVerificationGas, userOp.GasFees, conf.GetChainId(strategy.GetNewWork()), strategy.GetPaymasterAddress())
	if err != nil {
		return nil, "", err
	}
	userOpHash := crypto.Keccak256(byteRes)
	afterProcessUserOphash := utils.ToEthSignedMessageHash(userOpHash)
	return afterProcessUserOphash, hex.EncodeToString(byteRes), nil
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
