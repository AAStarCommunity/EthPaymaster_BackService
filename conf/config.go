package conf

import (
	contract_erc20 "AAStarCommunity/EthPaymaster_BackService/common/contract/erc20"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/tokens"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
)

var BasicConfig Config
var TokenContractCache map[common.Address]contract_erc20.Contract

func init() {
	BasicConfig = Config{}
}

type Config struct {
	NetworkConfigMap  map[network.Network]*NetWorkConfig `json:"network_config"`
	SupportEntryPoint map[network.Network]*mapset.Set[string]
	SupportPaymaster  map[network.Network]*mapset.Set[string]
}
type NetWorkConfig struct {
	ChainId     string                               `json:"chain_id"`
	IsTest      bool                                 `json:"is_test"`
	RpcUrl      string                               `json:"rpc_url"`
	ApiKey      string                               `json:"api_key"`
	TokenConfig map[tokens.TokenType]*common.Address `json:"token_config"`
}

func GetTokenAddress(networkParam network.Network, tokenParam tokens.TokenType) *common.Address {
	networkConfig := BasicConfig.NetworkConfigMap[networkParam]
	return networkConfig.TokenConfig[tokenParam]
}
func CheckEntryPointExist(network2 network.Network, address string) {

}
