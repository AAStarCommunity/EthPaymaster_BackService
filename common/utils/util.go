package utils

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"regexp"
	"strconv"
)

var HexPattern = regexp.MustCompile(`^0x[a-fA-F\d]*$`)

type EthCallReq struct {
	From common.Address `json:"from"`
	To   common.Address `json:"to"`
	Data hexutil.Bytes  `json:"data"`
}

type TraceCallOpts struct {
	Tracer         string      `json:"tracer"`
	StateOverrides OverrideSet `json:"stateOverrides"`
}
type OverrideSet map[common.Address]OverrideAccount
type OverrideAccount struct {
	Nonce     *hexutil.Uint64              `json:"nonce"`
	Code      *hexutil.Bytes               `json:"code"`
	Balance   *hexutil.Big                 `json:"balance"`
	State     *map[common.Hash]common.Hash `json:"state"`
	StateDiff *map[common.Hash]common.Hash `json:"stateDiff"`
}

func GenerateMockUservOperation() *map[string]any {
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

func PackIntTo32Bytes(left *big.Int, right *big.Int) [32]byte {
	leftBytes := left.Bytes()
	rightBytes := right.Bytes()

	leftHex := fmt.Sprintf("%016x", leftBytes)
	rightHex := fmt.Sprintf("%016x", rightBytes)

	leftBytes, _ = hex.DecodeString(leftHex)
	rightBytes, _ = hex.DecodeString(rightHex)

	var result [32]byte
	copy(result[:16], leftBytes)
	copy(result[16:], rightBytes)

	return result
}

func EncodeToStringWithPrefix(data []byte) string {
	res := hex.EncodeToString(data)
	if res[:2] != "0x" {
		return "0x" + res
	}
	return res
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
func SupplyZero(prefix string, maxTo int) string {
	padding := maxTo - len(prefix)
	if padding > 0 {
		prefix = "0" + prefix
		prefix = fmt.Sprintf("%0*s", maxTo, prefix)
	}
	return prefix
}
func IsLessThanZero(value *big.Int) bool {
	return false
	//TODO
}
func LeftIsLessTanRight(a *big.Int, b *big.Int) bool {
	return a.Cmp(b) < 0
}
func GetSign(message []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	digest := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(digest))
	sig, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}
	//In Ethereum, the last byte of the signature result represents the recovery ID of the signature, and by adding 27 to ensure it conforms to Ethereum's specification.
	sig[64] += 27
	return sig, nil
}
