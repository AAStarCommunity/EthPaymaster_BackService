package conf

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"fmt"
	"testing"
)

func TestConfigInit(t *testing.T) {
	BasicStrategyInit("../conf/basic_strategy_dev_config.json")
	BusinessConfigInit("../conf/business_dev_config.json")
	strategy := GetBasicStrategyConfig("Ethereum_Sepolia_v06_verifyPaymaster")
	if strategy == nil {
		t.Error("strategy is nil")
	}
	strategySuit, err := GetSuitableStrategy("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789", types.EthereumSepolia, types.PayTypeVerifying)
	if err != nil {
		t.Error("strategySuit is nil")
	}
	if strategySuit == nil {
		t.Error("strategySuit is nil")
	}

	chainid := GetChainId(types.EthereumSepolia)
	if chainid == "" {
		t.Error("chainid is 0")
	}
	fmt.Println(chainid)
	t.Log(chainid)

}
