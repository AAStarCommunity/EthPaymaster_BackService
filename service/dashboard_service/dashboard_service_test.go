package dashboard_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"encoding/json"
	"testing"
)

func TestDashBoardService(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	//t.Run("TestGetSuitableStrategy", TestGetSuitableStrategy)
	config.InitConfig("../../config/basic_strategy_config.json", "../../config/basic_config.json", "../../config/secret_config.json")
	Init()
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"TestGetSuitableStrategy",
			func(t *testing.T) {
				testGetStrategyByCode(t)
			},
		},
		{
			"TestGetSuitableStrategy",
			func(t *testing.T) {
				testGetSuitableStrategy(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
func testGetSuitableStrategy(t *testing.T) {
	strategy, err := GetSuitableStrategy("", global_const.EthereumSepolia, global_const.TokenTypeOP)
	if err != nil {
		t.Error(err)
	}
	jsonByte, _ := json.Marshal(strategy)
	t.Logf("Strategy: %s", string(jsonByte))
}
func testGetStrategyByCode(t *testing.T) {
	strategy, err := GetStrategyByCode("basic_arb_strategy__vHUZk", "", global_const.ArbitrumSpeolia)
	if err != nil {
		t.Error(err)
	}
	jsonByte, _ := json.Marshal(strategy)

	t.Logf("Strategy: %s", string(jsonByte))
}
