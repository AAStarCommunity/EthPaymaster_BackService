package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestOperator(t *testing.T) {
	conf.BasicStrategyInit("../../conf/basic_strategy_dev_config.json")
	conf.BusinessConfigInit("../../conf/business_dev_config.json")
	logrus.SetLevel(logrus.DebugLevel)
	mockRequest := getMockTryPayUserOpRequest()
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"testStrategyGenerate",
			func(t *testing.T) {
				testStrategyGenerate(t, mockRequest)
			},
		},
		{
			"testEstimateUserOpGas",
			func(t *testing.T) {
				testGetEstimateUserOpGas(t, mockRequest)
			},
		},
		{
			"testGetSupportEntrypointExecute",
			func(t *testing.T) {
				testGetSupportEntrypointExecute(t)
			},
		},
		{
			"Test_EthereumSepoliaV06Verify_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest.ForceStrategyId = string(global_const.StrategyCodeEthereumSepoliaV06Verify)
				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_OptimismSepoliaV06Verify_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest.ForceStrategyId = string(global_const.StrategyCodeOptimismSepoliaV06Verify)
				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_ArbitrumSpeoliaV06Verify_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest.ForceStrategyId = string(global_const.StrategyCodeArbitrumSpeoliaV06Verify)
				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_BaseSepoliaV06Verify_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest.ForceStrategyId = string(global_const.StrategyCodeArbitrumSpeoliaV06Verify)
				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_EthereumSepoliaV06Erc20_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest.Erc20Token = global_const.TokenTypeUSDT
				mockRequest.ForceStrategyId = string(global_const.StrategyCodeEthereumSepoliaV06Erc20)
				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_OpSepoliaV06Erc20_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest.Erc20Token = global_const.TokenTypeUSDT
				mockRequest.ForceStrategyId = string(global_const.StrategyCodeOptimismSepoliaV06Erc20)
				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_ArbSepoliaV06Erc20_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest.Erc20Token = global_const.TokenTypeUSDT
				mockRequest.ForceStrategyId = string(global_const.StrategyCodeArbitrumSpeoliaV06Erc20)
				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_BaseSepoliaV06Erc20_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest.Erc20Token = global_const.TokenTypeUSDT
				mockRequest.ForceStrategyId = string(global_const.StrategyCodeBaseSepoliaV06Erc20)
				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"testGetSupportStrategyExecute",
			func(t *testing.T) {
				testGetSupportStrategyExecute(t)
			},
		},
	}
	for _, tt := range tests {
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
	res, err := GetSupportEntrypointExecute("network")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}
func testTryPayUserOpExecute(t *testing.T, request *model.UserOpRequest) {
	result, err := TryPayUserOpExecute(request)
	if err != nil {
		t.Error(err)
		return
	}
	resultJson, _ := json.Marshal(result)
	t.Logf("Result: %v", string(resultJson))
}

func getMockTryPayUserOpRequest() *model.UserOpRequest {
	return &model.UserOpRequest{
		ForceStrategyId: "Ethereum_Sepolia_v06_verifyPaymaster",
		UserOp:          *utils.GenerateMockUservOperation(),
	}
}

func testGetSupportStrategyExecute(t *testing.T) {
	res, err := GetSupportStrategyExecute("network")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)

}
