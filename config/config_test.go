package config

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"testing"
)

func TestSecretConfigInit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	secretConfigInit("../config/secret_config.json")
	config := GetNetworkSecretConfig(global_const.EthereumSepolia)
	t.Log(config.RpcUrl)
	t.Log(config.SignerKey)
	t.Log(GetSigner(global_const.EthereumSepolia).Address.Hex())
}
func TestConfigInit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	InitConfig("../config/basic_strategy_config.json", "../config/basic_config.json", "../config/secret_config.json")
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
	rpcUrl := GetNewWorkClientURl(global_const.EthereumSepolia)
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
