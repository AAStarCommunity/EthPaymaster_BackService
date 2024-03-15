package utils

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	"regexp"
	"strconv"
)

var HexPattern = regexp.MustCompile(`^0x[a-fA-F\d]*$`)

func GenerateMockUserOperation() *map[string]any {
	//TODO use config
	var MockUserOpData = map[string]any{
		"call_data":                "0xb61d27f60000000000000000000000001c7d4b196cb0c7b01d743fbc6116a902379c7238000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000044095ea7b30000000000000000000000000000000000325602a77416a16136fdafd04b299fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000000000000000000000000000000000000000000000000000",
		"call_gas_limit":           "0x54fa",
		"init_code":                "0x9406cc6185a346906296840746125a0e449764545fbfb9cf000000000000000000000000340966abb6e37a06014546e0542b3aafad4550810000000000000000000000000000000000000000000000000000000000000000",
		"max_fee_per_gas":          "0x2aa887baca",
		"max_priority_fee_per_gas": "0x59682f00",
		"nonce":                    "0x00",
		"pre_verification_gas":     "0xae64",
		"sender":                   "0xF8498599744BC37e141cb800B67Dbf103a6b5881",
		"signature":                "0xfffffffffffffffffffffffffffffff0000000000000000000000000000000007aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1c",
		"verification_gas_limit":   "0x05fa35",
		"paymaster_and_data":       "0x0000000000325602a77416A16136FDafd04b299f",
	}

	return &MockUserOpData
}
func ValidateHex(value string) bool {
	if HexPattern.MatchString(value) {
		return true
	}
	return false
}
func IsStringInUint64Range(s string) bool {
	num, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return false
	}
	// 0 <= num <= MaxUint64
	return num <= ^uint64(0)
}
func GenerateUserOperation() *model.UserOperation {
	return &model.UserOperation{}
}
func EncodeToStringWithPrefix(data []byte) string {
	res := hex.EncodeToString(data)
	if res[:2] != "0x" {
		return "0x" + res
	}
	return res
}

func SignUserOp(privateKeyHex string, userOp *model.UserOperation) ([]byte, error) {

	serializedUserOp, err := json.Marshal(userOp)
	if err != nil {
		return nil, err
	}

	signature, err := SignMessage(privateKeyHex, string(serializedUserOp))
	if err != nil {
		return nil, err
	}

	return signature, nil
}
func SignMessage(privateKeyHex string, message string) ([]byte, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	// Hash the message
	hash := crypto.Keccak256([]byte(message))

	// Sign the hash
	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return nil, err
	}

	return signature, nil
}
