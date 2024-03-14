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
		"sender":                   "0x4A2FD3215420376DA4eD32853C19E4755deeC4D1",
		"nonce":                    "1",
		"init_code":                "0xe19e9755942bb0bd0cccce25b1742596b8a8250b3bf2c3e700000000000000000000000078d4f01f56b982a3b03c4e127a5d3afa8ebee6860000000000000000000000008b388a082f370d8ac2e2b3997e9151168bd09ff50000000000000000000000000000000000000000000000000000000000000000",
		"call_data":                "0xb61d27f6000000000000000000000000c206b552ab127608c3f666156c8e03a8471c72df000000000000000000000000000000000000000000000000002386f26fc1000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000000",
		"call_gas_limit":           "39837",
		"verification_gas_limit":   "100000",
		"max_fee_per_gas":          "44020",
		"max_priority_fee_per_gas": "1743509478",
		"pre_verification_gas":     "44020",
		"signature":                "0x760868cd7d9539c6e31c2169c4cab6817beb8247516a90e4301e929011451658623455035b83d38e987ef2e57558695040a25219c39eaa0e31a0ead16a5c925c1c",
		"paymaster_and_data":       "0x",
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
