package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strconv"
	"testing"
)

func TestTryPayUserOpExecute(t *testing.T) {
	request := getMockTryPayUserOpRequest()
	result, err := TryPayUserOpExecute(request)
	assert.NoError(t, err)
	resultJson, _ := json.Marshal(result)
	fmt.Printf("Result: %v", string(resultJson))
}

func getMockTryPayUserOpRequest() *model.TryPayUserOpRequest {
	return &model.TryPayUserOpRequest{
		ForceStrategyId: "1",
		UserOp:          *utils.GenerateMockUserOperation(),
	}
}

func TestGenerateTestData(t *testing.T) {
	strategy := dashboard_service.GetStrategyById("1")
	userOp, _ := userop.NewUserOp(utils.GenerateMockUserOperation())
	str, signature, err := generatePayMasterAndData(userOp, strategy)
	assert.NoError(t, err)
	fmt.Println(str)
	fmt.Println(signature)
	fmt.Println(len(signature))
}
func TestPackUserOp(t *testing.T) {
	// give same len signuature and paymasteranddata
	userOp, _ := userop.NewUserOp(utils.GenerateMockUserOperation())
	res, byteres, err := packUserOp(userOp)
	shouldEqualStr := "000000000000000000000000ffdb071c2b58ccc10ad386f9bb4e8d3d664ce73c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000000054fa000000000000000000000000000000000000000000000000000000000005fa35000000000000000000000000000000000000000000000000000000000000ae640000000000000000000000000000000000000000000000000000000059682f8e0000000000000000000000000000000000000000000000000000000059682f00000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000003c000000000000000000000000000000000000000000000000000000000000000589406cc6185a346906296840746125a0e449764545fbfb9cf000000000000000000000000b6bcf9517d193f551d0e3d6860103972dd13de7b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e4b61d27f60000000000000000000000001c7d4b196cb0c7b01d743fbc6116a902379c7238000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000044095ea7b30000000000000000000000000000000000325602a77416a16136fdafd04b299fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	assert.NoError(t, err)
	assert.EqualValues(t, shouldEqualStr, res)
	fmt.Println(res)
	fmt.Println(shouldEqualStr)
	fmt.Println(byteres)
}
func TestConvertHex(t *testing.T) {
	hexString := strconv.FormatUint(1500000000, 16)
	fmt.Println(hexString)
}
func TestGetValidTime(t *testing.T) {
	validStart, validEnd := getValidTime()
	fmt.Printf("validStart: %s, validEnd: %s\n", validStart, validEnd)
}
func TestSignPaymaster(t *testing.T) {

	userOp, _ := userop.NewUserOp(utils.GenerateMockUserOperation())
	strategy := dashboard_service.GetStrategyById("1")

	validStart, validEnd := getValidTime()

	//fmt.Printf("validStart: %s, validEnd: %s\n", validStart, validEnd)
	//message := fmt.Sprintf("%s%s%s%s", strategy.PayMasterAddress, string(strategy.PayType), validStart, validEnd)
	signatureByte, hashByte, err := SignPaymaster(userOp, strategy, validStart, validEnd)
	//signatureStr := hex.EncodeToString(signatureByte)
	assert.NoError(t, err)

	signatureStr := hex.EncodeToString(signatureByte)
	hashByteStr := hex.EncodeToString(hashByte)
	fmt.Printf("signatureStr len: %s\n", signatureStr)
	fmt.Printf("hashByteStr len: %s\n", hashByteStr)
	//
	//sigPublicKey, err := crypto.Ecrecover(hashByte, signatureByte)
	//if err != nil {
	//	assert.NoError(t, err)
	//	return
	//}
	//sigPublicKeyStr := hex.EncodeToString(sigPublicKey)
	//fmt.Println(sigPublicKeyStr)
	//publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	//if !ok {
	//	assert.Error(t, err)
	//	return
	//}
	//address := crypto.PubkeyToAddress(*publicKeyECDSA)
	//fmt.Printf("address: %s\n", address.Hex())
	//sigPublicKey, err := crypto.Ecrecover(hashToEthSignHash, signatureByte)
	//fmt.Println(sigPublicKey)
	//sigPublicKeyStr := hex.EncodeToString(sigPublicKey)
	//fmt.Println(sigPublicKeyStr)
	//assert.Equal(t, sigPublicKeyStr, "044eaed6b1f16e60354156fa334a094affc76d7b7061875a0b04290af9a14cc14ce2bce6ceba941856bd55c63f8199f408fff6495ce9d4c76899055972d23bdb3e")
	//
	//pubKeyBytes, _ := hex.DecodeString("044eaed6b1f16e60354156fa334a094affc76d7b7061875a0b04290af9a14cc14ce2bce6ceba941856bd55c63f8199f408fff6495ce9d4c76899055972d23bdb3e")
	//signatureNoRecoverID := signatureByte[:len(signatureByte)-1]
	//verified := crypto.VerifySignature(pubKeyBytes, hashToEthSignHash, signatureNoRecoverID)
	//assert.True(t, verified)

}

func TestSign(t *testing.T) {
	//hash 3244304e46b095a6dc5ff8af5cac03cbb22f6e07d3a0841dc4b3b8bc399a44702724cc7aad26b3854545269e34c156565f717b96acc52ee9de95526c644ddf6d00
	//sign  9429db04bd812b79bf15d55ee271426894cbfb6e7431da8d934d5e970dbf992c
	// address
}

func TestUserOpHash(t *testing.T) {
	strategy := dashboard_service.GetStrategyById("1")
	userOp, _ := userop.NewUserOp(utils.GenerateMockUserOperation())
	encodeHash, userOpabiEncodeStr, err := UserOpHash(userOp, strategy, big.NewInt(1710044496), big.NewInt(1741580496))
	assert.NoError(t, err)
	shouldEqualStr := "00000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000aa36a7000000000000000000000000d93349ee959d295b115ee223af10ef432a8e852300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000065ed35500000000000000000000000000000000000000000000000000000000067ce68d00000000000000000000000000000000000000000000000000000000000000300000000000000000000000000ffdb071c2b58ccc10ad386f9bb4e8d3d664ce73c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000000054fa000000000000000000000000000000000000000000000000000000000005fa35000000000000000000000000000000000000000000000000000000000000ae64000000000000000000000000000000000000000000000000000000005968334e0000000000000000000000000000000000000000000000000000000059682f00000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000003c000000000000000000000000000000000000000000000000000000000000000589406cc6185a346906296840746125a0e449764545fbfb9cf000000000000000000000000b6bcf9517d193f551d0e3d6860103972dd13de7b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e4b61d27f60000000000000000000000001c7d4b196cb0c7b01d743fbc6116a902379c7238000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000044095ea7b30000000000000000000000000000000000325602a77416a16136fdafd04b299fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	fmt.Printf("userOpabiEncodeStr %s \n", userOpabiEncodeStr)
	fmt.Printf("encodeHash %s \n", hex.EncodeToString(encodeHash))

	fmt.Println(shouldEqualStr)

	assert.EqualValues(t, userOpabiEncodeStr, shouldEqualStr)
	if userOpabiEncodeStr != shouldEqualStr {
		return
	}
	//fmt.Println("finish euqal abiencode")
	//userOpHashStr := hex.EncodeToString(userOpHash)
	//fmt.Println(userOpHashStr)
	//shouldEqualHashStr := "8ad4946fb4665c29754b83495e796fa03013aaa0f194326afad73ce2fc5b91e9"
	//assert.EqualValues(t, userOpHashStr, shouldEqualHashStr)
	//fmt.Println(userOpHashStr)
	//fmt.Println(shouldEqualHashStr)
}
func TestKeccak256(t *testing.T) {
	str := "00000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000aa36a7000000000000000000000000d93349ee959d295b115ee223af10ef432a8e852300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000065ed35500000000000000000000000000000000000000000000000000000000067ce68d00000000000000000000000000000000000000000000000000000000000000300000000000000000000000000ffdb071c2b58ccc10ad386f9bb4e8d3d664ce73c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000000054fa000000000000000000000000000000000000000000000000000000000005fa35000000000000000000000000000000000000000000000000000000000000ae64000000000000000000000000000000000000000000000000000000005968334e0000000000000000000000000000000000000000000000000000000059682f00000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000003c000000000000000000000000000000000000000000000000000000000000000589406cc6185a346906296840746125a0e449764545fbfb9cf000000000000000000000000b6bcf9517d193f551d0e3d6860103972dd13de7b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e4b61d27f60000000000000000000000001c7d4b196cb0c7b01d743fbc6116a902379c7238000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000044095ea7b30000000000000000000000000000000000325602a77416a16136fdafd04b299fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	//decimal, err := strconv.ParseInt(str, 16, 64)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(decimal)
	strByte, err := hex.DecodeString(str)
	if err != nil {
		fmt.Printf("has error  %s", err)
		return
	}
	fmt.Println(strByte)
	//strConvert := hex.EncodeToString(strByte)
	//fmt.Println(strConvert)
	//fmt.Println(strConvert)
	res := crypto.Keccak256(strByte)
	fmt.Println(hex.EncodeToString(res))

	//resHash := crypto.Keccak256Hash(strByte)
	//fmt.Println(resHash.Hex())
	//msg := []byte("abc")
	//exp, _ := hex.DecodeString("4e03657aea45a94fc7d47ba826c8d667c0d1e6e33a64a036ec44f58fa12d6c45")
	//checkhash(t, "Sha3-256-array", func(in []byte) []byte { h :=cry; return h[:] }, msg, exp)
}
func checkhash(t *testing.T, name string, f func([]byte) []byte, msg, exp []byte) {
	sum := f(msg)
	if !bytes.Equal(exp, sum) {
		t.Fatalf("hash %s mismatch: want: %x have: %x", name, exp, sum)
	}
}

func TestUserOP(t *testing.T) {
	userOp, _ := userop.NewUserOp(utils.GenerateMockUserOperation())
	fmt.Println(userOp.Sender.String())
}
func TestGenerateTestPaymaterDataparse(t *testing.T) {
	//contractABI, err := abi.JSON([]byte(`[
	//	{
	//		"constant": false,
	//		"inputs": [
	//			{
	//				"name": "userOp",
	//				"type": "tuple"
	//			},
	//			{
	//				"name": "requiredPreFund",
	//				"type": "uint256"
	//			}
	//		],
	//		"name": "_validatePaymasterUserOp",
	//		"outputs": [
	//			{
	//				"name": "context",
	//				"type": "bytes"
	//			},
	//			{
	//				"name": "validationData",
	//				"type": "uint256"
	//			}
	//		],
	//		"payable": false,
	//		"stateMutability": "nonpayable",
	//		"type": "function"
	//	}
	//]`))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//str := "0x
}
func TestDemo(t *testing.T) {
	//strategy := dashboard_service.GetStrategyById("1")
	userOp, _ := userop.NewUserOp(utils.GenerateMockUserOperation())

	//str := "0x"
	//fmt.Println(len(str))
	//fmt.Println(str[:2])
	//fmt.Println(str[:2] !=
	bytesTy, err := abi.NewType("bytes", "", nil)
	//uint256Ty, err := abi.NewType("uint256", "", nil)
	if err != nil {
		fmt.Println(err)
	}
	uint256Ty, _ := abi.NewType("uint256", "", nil)
	if err != nil {
		fmt.Println(err)
	}
	addressTy, _ := abi.NewType("address", "", nil)
	arguments := abi.Arguments{
		{
			Type: bytesTy,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: addressTy,
		},
	}
	packUserOpStr, _, err := packUserOp(userOp)
	//Btypelen := len(packUserOpStrByte)
	//byteArray := [Btypelen]byte(packUserOpStrByte)
	strByte, _ := hex.DecodeString(packUserOpStr)
	bytesRes, err := arguments.Pack(strByte, big.NewInt(1), common.Address{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hex.EncodeToString(bytesRes))
}
