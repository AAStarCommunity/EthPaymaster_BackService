package utils

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"log"
	"math/big"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var HexPattern = regexp.MustCompile(`^0x[a-fA-F\d]*$`)

const defaultStackSize = 10000

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
		"callData":             "0xb61d27f60000000000000000000000001c7d4b196cb0c7b01d743fbc6116a902379c7238000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000044095ea7b30000000000000000000000000000000000325602a77416a16136fdafd04b299fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00000000000000000000000000000000000000000000000000000000",
		"callGasLimit":         "0x54fa",
		"initCode":             "0x9406cc6185a346906296840746125a0e449764545fbfb9cf000000000000000000000000b6bcf9517d193f551d0e3d6860103972dd13de7b0000000000000000000000000000000000000000000000000000000000000000",
		"maxFeePerGas":         "0x5968606e",
		"maxPriorityFeePerGas": "0x59682f00",
		"nonce":                "0x00",
		"preVerificationGas":   "0xae64",
		"sender":               "0xFfDB071C2b58CCC10Ad386f9Bb4E8d3d664CE73c",
		"signature":            "0xaa846693598194980f3bf50486be854704534c1622d0c2ee895a5a1ebe1508221909a27cc7971d9f522c8df13b9d8a6ee446d09ea7635f31c59d77d35d1281421c",
		"verificationGasLimit": "0x05fa35",
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
	leftBytes = append(make([]byte, 16-len(leftBytes)), leftBytes...)
	rightBytes = append(make([]byte, 16-len(rightBytes)), rightBytes...)
	var result = [32]byte{}
	copy(result[:16], leftBytes)
	copy(result[16:], rightBytes)

	return result
}
func ConvertHexToBigInt(hex string) *big.Int {
	bigInt := new(big.Int)
	bigInt.SetString(hex, 0)
	return bigInt

}
func GetGasEntryPointGasPrice(maxFeePerGas *big.Int, maxPriorityFeePerGas *big.Int, baseFee *big.Int) *big.Int {
	if maxFeePerGas == maxPriorityFeePerGas {
		// is 1559 not support
		return maxFeePerGas
	}
	combineFee := new(big.Int).Add(baseFee, maxPriorityFeePerGas)
	return GetMinValue(maxFeePerGas, combineFee)
}

func EncodeToHexStringWithPrefix(data []byte) string {
	res := hex.EncodeToString(data)
	if res[:2] != "0x" {
		return "0x" + res
	}
	return res
}
func ConvertBigIntToHexWithPrefix(data *big.Int) string {
	res := fmt.Sprintf("0x%x", data)
	return res
}

func DecodeStringWithPrefix(data string) ([]byte, error) {
	if data[:2] == "0x" {
		data = data[2:]
	}
	return hex.DecodeString(data)
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

func GetMinValue(int2 *big.Int, int3 *big.Int) *big.Int {
	if int2.Cmp(int3) == -1 {
		return int2
	}
	return int3
}
func ConvertBalanceToEther(balance *big.Int) *big.Float {
	balanceFloat := new(big.Float).SetInt(balance)
	balanceFloat = new(big.Float).Quo(balanceFloat, global_const.EthWeiFactor)
	return balanceFloat
}
func CounverEtherToWei(ether *big.Float) *big.Int {
	afterEther := ether.Mul(ether, global_const.EthWeiFactor)
	afterEtherInt, _ := afterEther.Int(nil)
	return afterEtherInt
}
func ConvertStringToSet(input string, split string) mapset.Set[string] {
	set := mapset.NewSet[string]()
	arr := strings.Split(input, split)
	for _, value := range arr {
		set.Add(value)
	}
	return set
}
func ParseCallError(err error, abi *abi.ABI) (string, error) {
	rpcErr, ok := err.(rpc.DataError)
	if !ok {
		return "", xerrors.Errorf("ExecutionResult: cannot assert type: error is not of type rpc.DataError")
	}
	data, ok := rpcErr.ErrorData().(string)
	if !ok {
		return "", xerrors.Errorf("ExecutionResult: cannot assert type: data is not of type string")
	}
	logrus.Debugf("data :[%s]", data)

	for _, abiErr := range abi.Errors {
		logrus.Debugf("abiErr :[%v]", abiErr)
		revert, uppackErr := abiErr.Unpack(common.Hex2Bytes(data[2:]))
		if uppackErr != nil {
			logrus.Debugf("executionResult err: [%s]", uppackErr.Error())
		} else {
			logrus.Debugf("has revert :[%v]", revert)

			return fmt.Sprintf("%v", revert), nil
		}
	}
	return "", nil
}

func GetCurrentGoroutineStack() string {
	var buf [defaultStackSize]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}
func DBTransactional(db *gorm.DB, handle func(*gorm.DB) error) (err error) {
	tx := db.Begin()
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			logrus.Errorf("TX ERROR [%s] ", GetCurrentGoroutineStack())
			err = xerrors.Errorf("TX ERROR [%v]", p)
			//panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()
	err = handle(tx)
	return err
}

func TransEth(from *ecdsa.PrivateKey, toAddress *common.Address, client *ethclient.Client, amount *big.Int, chainID *big.Int) (*types.Transaction, error) {
	fromPrivateKey := from.Public()
	fromPublicKeyECDSA, ok := fromPrivateKey.(*ecdsa.PublicKey)
	if !ok {
		logrus.Error("error casting public key to ECDSA")
		return nil, xerrors.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*fromPublicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err

	}
	gasLimit := uint64(21000) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	var data []byte
	tx := types.NewTx(&types.DynamicFeeTx{

		Nonce:     nonce,
		Data:      data,
		Gas:       gasLimit,
		GasFeeCap: gasPrice,
		GasTipCap: gasPrice,
		Value:     amount,
		To:        toAddress,
	})

	signTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), from)

	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signTx)
	if err != nil {
		log.Fatal(err)
	}
	return signTx, nil
}
