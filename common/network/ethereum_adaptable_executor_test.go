package network

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"context"
	"testing"
)

func TestSimulateV06HandleOp(t *testing.T) {
	sepoliaExector := GetEthereumExecutor(types.EthereumSepolia)
	op, newErr := userop.NewUserOp(utils.GenerateMockUserOperation(), types.EntrypointV06)
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
func TestEthereumExecutorClientConnect(t *testing.T) {
	conf.BasicStrategyInit()
	conf.BusinessConfigInit()
	executor := GetEthereumExecutor(types.EthereumSepolia)
	client := executor.Client
	client.ChainID(context.Background())
}
