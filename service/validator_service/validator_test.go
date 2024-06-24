package validator_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"testing"
)

func TestValidatorService(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.InitConfig("../../config/basic_strategy_config.json", "../../config/basic_config.json", "../../config/secret_config.json")
	dashboard_service.Init()
	strategyCode := "basic_arb_strategy__vHUZk"
	strategy, err := dashboard_service.GetStrategyByCode(strategyCode, "", global_const.ArbitrumSpeolia)
	request := getMockTryPayUserOpRequest(strategyCode, global_const.ArbitrumSpeolia)
	if err != nil {
		t.Fatalf("GetStrategyByCode error: %v", err)
	}
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"TestValidateStrategy",
			func(t *testing.T) {
				testValidateStrategy(t, strategy, request)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func testValidateStrategy(t *testing.T, strategy *model.Strategy, request *model.UserOpRequest) {
	if err := ValidateStrategy(strategy, request, &model.ApiKeyModel{
		UserId: 5,
	}); err != nil {
		t.Fatalf("ValidateStrategy error: %v", err)
	}
}

func getMockTryPayUserOpRequest(strategyCode string, network global_const.Network) *model.UserOpRequest {
	return &model.UserOpRequest{
		StrategyCode: strategyCode,
		Network:      network,
		UserOp:       *utils.GenerateMockUservOperation(),
	}
}
