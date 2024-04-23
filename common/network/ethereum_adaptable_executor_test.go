package network

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"context"
	"github.com/stretchr/testify/assert"
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
			"TestOptimismSepoliaClientConnect",
			func(t *testing.T) {
				testEthereumExecutorClientConnect(t, types.OptimismSepolia)
			},
		},
		{
			"TestScrollSepoliaClientConnect",
			func(t *testing.T) {
				testEthereumExecutorClientConnect(t, types.ScrollSepolia)
			},
		},
		{
			"TestArbitrumSpeoliaClientConnect",
			func(t *testing.T) {
				testEthereumExecutorClientConnect(t, types.ArbitrumSpeolia)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
func TestSimulateV06HandleOp(t *testing.T) {
	sepoliaExector := GetEthereumExecutor(types.EthereumSepolia)
	op, newErr := userop.NewUserOp(utils.GenerateMockUserv06Operation(), types.EntrypointV06)
	if newErr != nil {
		return
	}
	opValue := *op
	strategy := conf.GetBasicStrategyConfig("Ethereum_Sepolia_v06_verifyPaymaster")
	userOpV6 := opValue.(*userop.UserOperationV06)
	simulataResult, err := sepoliaExector.SimulateV06HandleOp(userOpV6, strategy.GetEntryPointAddress())
	if err != nil {
		return
	}
	if simulataResult == nil {
		return
	}
}

func testEthereumExecutorClientConnect(t *testing.T, chain types.Network) {
	executor := GetEthereumExecutor(chain)
	client := executor.Client
	chainId, err := client.ChainID(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, chainId)
	assert.Equal(t, chainId, executor.ChainId)
}
