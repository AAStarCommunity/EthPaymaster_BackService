package arbitrum

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"encoding/json"
	"testing"
)

func TestGetArbitrumGas(t *testing.T) {
	config.InitConfig("../../../config/basic_strategy_config.json", "../../../config/basic_config.json", "../../../config/secret_config.json")
	strategy := config.GetBasicStrategyConfig("Arbitrum_Sepolia_v06_verifyPaymaster")
	if strategy == nil {
		t.Error("strategy is nil")
	}
	op, err := user_op.NewUserOp(utils.GenerateMockUservOperation())
	if err != nil {
		t.Error(err)
		return
	}
	executor := network.GetEthereumExecutor(strategy.GetNewWork())
	gasOutPut, err := GetArbEstimateOutPut(executor.Client, &model.PreVerificationGasEstimateInput{
		Strategy:         strategy,
		Op:               op,
		GasFeeResult:     &model.GasPrice{},
		SimulateOpResult: model.MockSimulateHandleOpResult,
	})
	if err != nil {
		t.Error(err)
		return
	}
	jsonRes, _ := json.Marshal(gasOutPut)
	t.Logf("gasOutPut:%v", string(jsonRes))
}
