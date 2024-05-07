package config

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"testing"
)

func TestConfigInit(t *testing.T) {
	BasicStrategyInit("../config/basic_strategy_dev_config.json")
	BusinessConfigInit("../config/business_dev_config.json")
	strategy := GetBasicStrategyConfig("Ethereum_Sepolia_v06_verifyPaymaster")
	if strategy == nil {
		t.Error("strategy is nil")
		return
	}
	strategySuit, err := GetSuitableStrategy(global_const.EntrypointV06, global_const.EthereumSepolia, global_const.PayTypeVerifying)
	if err != nil {

		t.Error("strategySuit is nil")
		return
	}
	if strategySuit == nil {
		t.Error("strategySuit is nil")
		return
	}

	chainId := GetChainId(global_const.EthereumSepolia)
	if chainId == "" {
		t.Error("chainid is 0")
	}
	t.Log(chainId)
	rpcUrl := GetEthereumRpcUrl(global_const.EthereumSepolia)
	if rpcUrl == "" {
		t.Error("rpcUrl is 0")
	}
	t.Log(rpcUrl)

	eoa := GetSigner(global_const.EthereumSepolia)
	if eoa == nil {
		t.Error("eoa is nil")
	}
	t.Log(eoa.Address.Hex())
	scrollEoa := GetSigner(global_const.ScrollSepolia)
	if scrollEoa == nil {
		t.Error("eoa is nil")
	}
	t.Log(scrollEoa.Address.Hex())
}
