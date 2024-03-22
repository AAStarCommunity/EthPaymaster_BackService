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
		"max_fee_per_gas":          "0x5968606e",
		"max_priority_fee_per_gas": "0x59682f00",
		"nonce":                    "0x00",
		"pre_verification_gas":     "0xae64",
		"sender":                   "0xffdb071c2b58ccc10ad386f9bb4e8d3d664ce73c",
		"signature":                "0xaa846693598194980f3bf50486be854704534c1622d0c2ee895a5a1ebe1508221909a27cc7971d9f522c8df13b9d8a6ee446d09ea7635f31c59d77d35d1281421c",
		"verification_gas_limit":   "0x05fa35",
		"paymaster_and_data":       "0xd93349Ee959d295B115Ee223aF10EF432A8E8523000000000000000000000000000000000000000000000000000000000065ed35500000000000000000000000000000000000000000000000000000000067ce68d015fdcf36211b7269133323a60a4b783a6a91ff72f1c5ad31398e259b9be5bb980d1a07b3aaee9a1c3f4bcc37c64bbf3e86da1b30227ca7d737b940caef5778191b",
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
