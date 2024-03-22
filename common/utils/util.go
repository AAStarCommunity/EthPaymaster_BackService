package utils

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"bytes"
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
		"init_code":                "0x9406cc6185a346906296840746125a0e449764545fbfb9cf000000000000000000000000b6bcf9517d193f551d0e3d6860103972dd13de7b0000000000000000000000000000000000000000000000000000000000000000",
		"max_fee_per_gas":          "0x5968334e",
		"max_priority_fee_per_gas": "0x59682f00",
		"nonce":                    "0x00",
		"pre_verification_gas":     "0xae64",
		"sender":                   "0xffdb071c2b58ccc10ad386f9bb4e8d3d664ce73c",
		"signature":                "0xe0a9eca60d705c9bfa6a91d794cf9c3e00058892f42f16ea55105086b23be1726457daf05b32290789d357b2ba042ce4564dba690d5e4c2211ca11c300de94d21c",
		"verification_gas_limit":   "0x05fa35",
		"paymaster_and_data":       "0xE99c4Db5E360B8c84bF3660393CB2A85c3029b4400000000000000000000000000000000000000000000000000000000171004449600000000000000000000000000000000000000000000000000000017415804969e46721fc1938ac427add8a9e0d5cba2be5b17ccda9b300d0d3eeaff1904dfc23e276abd1ba6e3e269ec6aa36fe6a2442c18d167b53d7f9f0d1b3ebe80b09a6200",
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

func ToEthSignedMessageHash(msg []byte) []byte {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte("\x19Ethereum Signed Message:\n32"))
	buffer.Write(msg)
	return crypto.Keccak256(buffer.Bytes())
}

func ReplaceLastTwoChars(str, replacement string) string {
	if len(str) < 2 {
		return str
	}
	return str[:len(str)-2] + replacement
}
