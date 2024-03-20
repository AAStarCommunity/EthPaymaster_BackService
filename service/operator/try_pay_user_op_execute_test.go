package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"math/big"
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
	userOp, _ := model.NewUserOp(utils.GenerateMockUserOperation())
	str, signature, err := generatePayMasterAndData(userOp, strategy)
	assert.NoError(t, err)
	fmt.Println(str)
	fmt.Println(signature)
	fmt.Println(len(signature))
}
func TestPackUserOp(t *testing.T) {
	// give same len signuature and paymasteranddata
	userOp, _ := model.NewUserOp(utils.GenerateMockUserOperation())
	res, byteres, err := packUserOp(userOp)
	shouldEqualStr := "000000000000000000000000f8498599744bc37e141cb800b67dbf103a6b58810000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000000054fa000000000000000000000000000000000000000000000000000000000005fa35000000000000000000000000000000000000000000000000000000000000ae640000000000000000000000000000000000000000000000000000002aa887baca0000000000000000000000000000000000000000000000000000000059682f00000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000003c000000000000000000000000000000000000000000000000000000000000000589406cc6185a346906296840746125a0e449764545fbfb9cf000000000000000000000000340966abb6e37a06014546e0542b3aafad4550810000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e4b61d27f60000000000000000000000001c7d4b196cb0c7b01d743fbc6116a902379c7238000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000044095ea7b30000000000000000000000000000000000325602a77416a16136fdafd04b299fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	assert.NoError(t, err)
	assert.EqualValues(t, shouldEqualStr, res)
	fmt.Println(res)
	fmt.Println(shouldEqualStr)
	fmt.Println(byteres)
}

func TestUserOpHash(t *testing.T) {
	strategy := dashboard_service.GetStrategyById("1")
	userOp, _ := model.NewUserOp(utils.GenerateMockUserOperation())
	userOpHash, userOpabiEncodeStr, err := UserOpHash(userOp, strategy, big.NewInt(1710044496), big.NewInt(1741580496))
	assert.NoError(t, err)
	shouldEqualStr := "00000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000aa36a7000000000000000000000000d93349ee959d295b115ee223af10ef432a8e852300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000065ed35500000000000000000000000000000000000000000000000000000000067ce68d00000000000000000000000000000000000000000000000000000000000000300000000000000000000000000f8498599744bc37e141cb800b67dbf103a6b58810000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000000054fa000000000000000000000000000000000000000000000000000000000005fa35000000000000000000000000000000000000000000000000000000000000ae640000000000000000000000000000000000000000000000000000002aa887baca0000000000000000000000000000000000000000000000000000000059682f00000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000003c000000000000000000000000000000000000000000000000000000000000000589406cc6185a346906296840746125a0e449764545fbfb9cf000000000000000000000000340966abb6e37a06014546e0542b3aafad4550810000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e4b61d27f60000000000000000000000001c7d4b196cb0c7b01d743fbc6116a902379c7238000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000044095ea7b30000000000000000000000000000000000325602a77416a16136fdafd04b299fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	fmt.Println(userOpabiEncodeStr)
	fmt.Println(shouldEqualStr)
	assert.EqualValues(t, userOpabiEncodeStr, shouldEqualStr)
	userOpHashStr := hex.EncodeToString(userOpHash)
	fmt.Println(userOpHashStr)
	shouldEqualHashStr := "a1e2c52ad5779f4eb7a87c570149d7d33614fbc1d1ac30fa6cfe80107909e0fa"
	assert.EqualValues(t, userOpHashStr, shouldEqualHashStr)
}

func TestUserOP(t *testing.T) {
	userOp, _ := model.NewUserOp(utils.GenerateMockUserOperation())
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
	userOp, _ := model.NewUserOp(utils.GenerateMockUserOperation())

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
