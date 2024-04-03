package conf

import (
	contract_erc20 "AAStarCommunity/EthPaymaster_BackService/common/contract/erc20"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"github.com/ethereum/go-ethereum/common"
	"go/token"
)

var BasicConfig Config
var TokenContractCache map[common.Address]contract_erc20.Contract

func init() {
	BasicConfig = Config{}
}

type Config struct {
	NetworkConfigMap  map[network.Network]*NetWorkConfig `json:"network_config"`
	SupportEntryPoint map[network.Network][]string
}
type NetWorkConfig struct {
	ChainId     string                          `json:"chain_id"`
	IsTest      bool                            `json:"is_test"`
	RpcUrl      string                          `json:"rpc_url"`
	ApiKey      string                          `json:"api_key"`
	TokenConfig map[token.Token]*common.Address `json:"token_config"`
}

func GetTokenAddress(networkParam network.Network, tokenParam token.Token) *common.Address {
	networkConfig := BasicConfig.NetworkConfigMap[networkParam]
	return networkConfig.TokenConfig[tokenParam]
}
