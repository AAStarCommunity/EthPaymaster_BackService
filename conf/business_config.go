package conf

import (
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_common/contract/erc20"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var BasicConfig BusinessConfig
var TokenContractCache map[common.Address]contract_erc20.Contract

func init() {
	BasicConfig = BusinessConfig{}
}

type BusinessConfig struct {
	NetworkConfigMap  map[types.Network]*NetWorkConfig `json:"network_config"`
	SupportEntryPoint map[types.Network]*mapset.Set[string]
	SupportPaymaster  map[types.Network]*mapset.Set[string]
}
type NetWorkConfig struct {
	ChainId     *big.Int                            `json:"chain_id"`
	IsTest      bool                                `json:"is_test"`
	RpcUrl      string                              `json:"rpc_url"`
	ApiKey      string                              `json:"api_key"`
	TokenConfig map[types.TokenType]*common.Address `json:"token_config"`
	GasToken    types.TokenType
}

func GetTokenAddress(networkParam types.Network, tokenParam types.TokenType) *common.Address {
	networkConfig := BasicConfig.NetworkConfigMap[networkParam]
	return networkConfig.TokenConfig[tokenParam]
}
func CheckEntryPointExist(network2 types.Network, address string) bool {
	entryPointSet := BasicConfig.SupportEntryPoint[network2]
	entryPointSetValue := *entryPointSet
	return entryPointSetValue.Contains(address)
}
func GetGasToken(networkParam types.Network) types.TokenType {
	networkConfig := BasicConfig.NetworkConfigMap[networkParam]
	return networkConfig.GasToken
}

func GetChainId(newworkParam types.Network) *big.Int {
	networkConfig := BasicConfig.NetworkConfigMap[newworkParam]
	return networkConfig.ChainId
}
