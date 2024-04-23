package network

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"context"
	"testing"
)

func TestEthereumAdaptableExecutor(t *testing.T) {
	conf.BasicStrategyInit("../../conf/basic_strategy_dev_config.json")
	conf.BusinessConfigInit("../../conf/business_dev_config.json")
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"TestEthereumSepoliaClientConnect",
			func(t *testing.T) {
				testEthereumExecutorClientConnect(t, types.EthereumSepolia)
			},
		},

		{
			"TestScrollSepoliaClientConnect",
			func(t *testing.T) {
				testEthereumExecutorClientConnect(t, types.ScrollSepolia)
			},
		},
		{
			"TestOptimismSepoliaClientConnect",
			func(t *testing.T) {
				testEthereumExecutorClientConnect(t, types.OptimismSepolia)
			},
		},
		{
			"TestArbitrumSpeoliaClientConnect",
			func(t *testing.T) {
				testEthereumExecutorClientConnect(t, types.ArbitrumSpeolia)
			},
		},
		{
			"TestSepoliaSimulateV06HandleOp",
			func(t *testing.T) {
				testSimulateV06HandleOp(t, types.ArbitrumSpeolia)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
func testSimulateV06HandleOp(t *testing.T, chain types.Network) {
	sepoliaExector := GetEthereumExecutor(chain)
	op, newErr := user_op.NewUserOp(utils.GenerateMockUservOperation(), types.EntrypointV06)
	if newErr != nil {
		t.Error(newErr)
	}
	strategy := conf.GetBasicStrategyConfig("Ethereum_Sepolia_v06_verifyPaymaster")
	simulataResult, err := sepoliaExector.SimulateV06HandleOp(op, strategy.GetEntryPointAddress())
	if err != nil {
		t.Error(err)
	}
	if simulataResult == nil {
		t.Error("simulataResult is nil")
	}
	t.Logf("simulateResult: %v", simulataResult)
}

func testEthereumExecutorClientConnect(t *testing.T, chain types.Network) {
	executor := GetEthereumExecutor(chain)
	if executor == nil {
		t.Error("executor is nil")
	}
	client := executor.Client
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		t.Error(err)
	}
	if chainId == nil {
		t.Error("chainId is nil")
	}
	if chainId.String() != executor.ChainId.String() {
		t.Errorf(" %s chainId not equal %s", chainId.String(), executor.ChainId.String())
	}
	t.Logf("network %s chainId: %s", chain, chainId.String())
}
