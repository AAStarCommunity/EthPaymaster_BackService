package network

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"context"
	"testing"
)

func SimulateV06HandleOp() {

}
func TestEthereumExecutorClientConnect(t *testing.T) {
	conf.BasicStrategyInit()
	conf.BusinessConfigInit()
	executor := GetEthereumExecutor(types.EthereumSepolia)
	client := executor.Client
	client.ChainID(context.Background())
}
