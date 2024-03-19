package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
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

	assert.NoError(t, err)
	fmt.Println(res)
	fmt.Println(byteres)
}

func TestUserOpHash(t *testing.T) {
	validStart, validEnd := getValidTime()
	strategy := dashboard_service.GetStrategyById("1")
	userOp, _ := model.NewUserOp(utils.GenerateMockUserOperation())
	userOpHash, err := UserOpHash(userOp, strategy, validStart, validEnd)
	assert.NoError(t, err)
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
