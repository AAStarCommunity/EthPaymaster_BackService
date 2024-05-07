package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestOperator(t *testing.T) {
	config.BasicStrategyInit("../../config/basic_strategy_dev_config.json")
	config.BusinessConfigInit("../../config/business_dev_config.json")
	logrus.SetLevel(logrus.DebugLevel)
	immutableRequest := getMockTryPayUserOpRequest()
	mockRequestNotSupport1559 := getMockTryPayUserOpRequest()
	mockRequestNotSupport1559.UserOp["maxPriorityFeePerGas"] = mockRequestNotSupport1559.UserOp["maxFeePerGas"]
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"testStrategyGenerate",
			func(t *testing.T) {
				testStrategyGenerate(t, immutableRequest)
			},
		},
		{
			"testEstimateUserOpGas",
			func(t *testing.T) {
				testGetEstimateUserOpGas(t, immutableRequest)
			},
		},
		{
			"testGetSupportEntrypointExecute",
			func(t *testing.T) {
				testGetSupportEntrypointExecute(t)
			},
		},
		{
			"Test_ScrollSepoliaV06Verify_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequestNotSupport1559.ForceStrategyId = string(global_const.StrategyCodeScrollSepoliaV06Verify)
				mockRequest := getMockTryPayUserOpRequest()

				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_EthereumSepoliaV06Verify_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest := getMockTryPayUserOpRequest()

				mockRequest.ForceStrategyId = string(global_const.StrategyCodeEthereumSepoliaV06Verify)

				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_OptimismSepoliaV06Verify_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest := getMockTryPayUserOpRequest()

				mockRequest.ForceStrategyId = string(global_const.StrategyCodeOptimismSepoliaV06Verify)

				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_ArbitrumSpeoliaV06Verify_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest := getMockTryPayUserOpRequest()

				mockRequest.ForceStrategyId = string(global_const.StrategyCodeArbitrumSepoliaV06Verify)

				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_BaseSepoliaV06Verify_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest := getMockTryPayUserOpRequest()
				mockRequest.ForceStrategyId = string(global_const.StrategyCodeArbitrumSepoliaV06Verify)

				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_EthereumSepoliaV06Erc20_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest := getMockTryPayUserOpRequest()

				mockRequest.Erc20Token = global_const.TokenTypeUSDT
				mockRequest.ForceStrategyId = string(global_const.StrategyCodeEthereumSepoliaV06Erc20)
				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_OpSepoliaV06Erc20_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest := getMockTryPayUserOpRequest()

				mockRequest.Erc20Token = global_const.TokenTypeUSDT
				mockRequest.ForceStrategyId = string(global_const.StrategyCodeOptimismSepoliaV06Erc20)
				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_ArbSepoliaV06Erc20_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest := getMockTryPayUserOpRequest()
				mockRequest.Erc20Token = global_const.TokenTypeUSDT
				mockRequest.ForceStrategyId = string(global_const.StrategyCodeArbitrumSepoliaV06Erc20)
				testTryPayUserOpExecute(t, mockRequest)
			},
		},
	}
	for _, tt := range tests {
		if os.Getenv("GITHUB_ACTIONS") != "" && global_const.GitHubActionWhiteListSet.Contains(tt.name) {
			t.Logf("Skip test [%s] in GitHub Actions", tt.name)
			continue
		}
		t.Run(tt.name, tt.test)
	}

}
func testGetEstimateUserOpGas(t *testing.T, request *model.UserOpRequest) {
	result, err := GetEstimateUserOpGas(request)
	if err != nil {
		t.Error(err)
		return
	}
	resultJson, _ := json.Marshal(result)
	fmt.Printf("Result: %v", string(resultJson))
}
func testStrategyGenerate(t *testing.T, request *model.UserOpRequest) {
	strategy, err := StrategyGenerate(request)
	if err != nil {
		t.Error(err)
		return
	}
	strategyJson, _ := json.Marshal(strategy)
	fmt.Printf("Strategy: %v", string(strategyJson))
}
func testGetSupportEntrypointExecute(t *testing.T) {
	res, err := GetSupportEntrypointExecute("ethereum-sepolia")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}
func testTryPayUserOpExecute(t *testing.T, request *model.UserOpRequest) {
	result, err := TryPayUserOpExecute(request)
	if err != nil {
		t.Fatal(err)
		return
	}
	resultJson, _ := json.Marshal(result)
	t.Logf("Result: %v", string(resultJson))

	executor := network.GetEthereumExecutor(result.NetWork)
	if executor == nil {
		t.Error("executor is nil")
		return
	}
	userOp, err := user_op.NewUserOp(&request.UserOp)
	paymasterDataStr := result.UserOpResponse.PayMasterAndData

	paymasterData, err := utils.DecodeStringWithPrefix(paymasterDataStr)
	userOp.PaymasterAndData = paymasterData
	if err != nil {
		t.Error(err)
		return
	}
	if result.EntrypointVersion == global_const.EntrypointV07 {
		//TODO
	} else {
		userOp.VerificationGasLimit = result.UserOpResponse.VerificationGasLimit
		userOp.PreVerificationGas = result.UserOpResponse.PreVerificationGas
		userOp.MaxFeePerGas = result.UserOpResponse.MaxFeePerGas
		userOp.MaxPriorityFeePerGas = result.UserOpResponse.MaxPriorityFeePerGas
		userOp.CallGasLimit = result.UserOpResponse.CallGasLimit
		address := common.HexToAddress(result.EntryPointAddress)
		jsonUserOP, err := json.Marshal(userOp)
		if err != nil {
			t.Error(err)
			return
		}
		t.Logf("jsonUserOP: %v", string(jsonUserOP))
		result, err := executor.SimulateV06HandleOp(userOp, &address)
		if err != nil {
			t.Fatal(err)
			return
		}
		resultJson, _ := json.Marshal(result)
		t.Logf("Result: %v", string(resultJson))
	}
}

func getMockTryPayUserOpRequest() *model.UserOpRequest {
	return &model.UserOpRequest{
		ForceStrategyId: "Ethereum_Sepolia_v06_verifyPaymaster",
		UserOp:          *utils.GenerateMockUservOperation(),
	}
}
