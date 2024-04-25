package conf

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"testing"
)

func TestConfigInit(t *testing.T) {
	BasicStrategyInit("../conf/basic_strategy_dev_config.json")
	BusinessConfigInit("../conf/business_dev_config.json")
	strategy := GetBasicStrategyConfig("Ethereum_Sepolia_v06_verifyPaymaster")
	if strategy == nil {
		t.Error("strategy is nil")
	}
	strategySuit, err := GetSuitableStrategy("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789", global_const.EthereumSepolia, global_const.PayTypeVerifying)
	if err != nil {
		t.Error("strategySuit is nil")
	}
	if strategySuit == nil {
		t.Error("strategySuit is nil")
	}

	chainid := GetChainId(global_const.EthereumSepolia)
	if chainid == "" {
		t.Error("chainid is 0")
	}
	t.Log(chainid)
	rpcUrl := GetEthereumRpcUrl(global_const.EthereumSepolia)
	if rpcUrl == "" {
		t.Error("rpcUrl is 0")
	}
	t.Log(rpcUrl)

}
