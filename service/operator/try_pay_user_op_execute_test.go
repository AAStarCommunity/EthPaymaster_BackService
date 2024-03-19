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
	userOp, _ := model.NewUserOp(utils.GenerateMockUserOperation())
	userOp.PaymasterAndData = nil
	userOp.Signature = nil
	res, byteres, err := packUserOp(userOp)
	shouldEqualStr := "000000000000000000000000f8498599744bc37e141cb800b67dbf103a6b58810000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000000054fa000000000000000000000000000000000000000000000000000000000005fa35000000000000000000000000000000000000000000000000000000000000ae640000000000000000000000000000000000000000000000000000002aa887baca0000000000000000000000000000000000000000000000000000000059682f000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000034000000000000000000000000000000000000000000000000000000000000000589406cc6185a346906296840746125a0e449764545fbfb9cf000000000000000000000000340966abb6e37a06014546e0542b3aafad4550810000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e4b61d27f60000000000000000000000001c7d4b196cb0c7b01d743fbc6116a902379c7238000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000044095ea7b30000000000000000000000000000000000325602a77416a16136fdafd04b299fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	assert.NoError(t, err)
	if shouldEqualStr != res {
		fmt.Println("not equal")
	}
	fmt.Println(res)
	fmt.Println(shouldEqualStr)
	fmt.Println(byteres)
}

func TestUserOpHash(t *testing.T) {
	strategy := dashboard_service.GetStrategyById("1")
	userOp, _ := model.NewUserOp(utils.GenerateMockUserOperation())
	userOpHash, userOpHashStr, err := UserOpHash(userOp, strategy, big.NewInt(0xffffffffff), big.NewInt(0xaa))
	assert.NoError(t, err)
	shouldEqualStr := "00000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000aa36a7000000000000000000000000e99c4db5e360b8c84bf3660393cb2a85c3029b440000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000ffffffffff00000000000000000000000000000000000000000000000000000000000000aa0000000000000000000000000000000000000000000000000000000000000300000000000000000000000000f8498599744bc37e141cb800b67dbf103a6b58810000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000001e000000000000000000000000000000000000000000000000000000000000054fa000000000000000000000000000000000000000000000000000000000005fa35000000000000000000000000000000000000000000000000000000000000ae640000000000000000000000000000000000000000000000000000002aa887baca0000000000000000000000000000000000000000000000000000000059682f000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000034000000000000000000000000000000000000000000000000000000000000000589406cc6185a346906296840746125a0e449764545fbfb9cf000000000000000000000000340966abb6e37a06014546e0542b3aafad4550810000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e4b61d27f60000000000000000000000001c7d4b196cb0c7b01d743fbc6116a902379c7238000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000044095ea7b30000000000000000000000000000000000325602a77416a16136fdafd04b299fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	fmt.Println(userOpHashStr)
	fmt.Println(shouldEqualStr)
	fmt.Println(hex.EncodeToString(userOpHash))
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
