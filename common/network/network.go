package network

import "AAStarCommunity/EthPaymaster_BackService/common/erc20_token"

type NetworkInfo struct {
	Name     string                `json:"main_net_name"`
	RpcUrl   string                `json:"main_net_rpc_url"`
	GasToken erc20_token.TokenType `json:"gas_token"`
}

//newworkConfig : chainId,GasToken, name, is_test,
//newwork clinetconfig : name, rpc_url, apikey,

//type Chain string
//
//const (
//	Ethereum Chain = "Ethereum"
//	Arbitrum Chain = "Arbitrum"
//	Optimism Chain = "Optimism"
//)

type Network string

const (
	Ethereum Network = "ethereum"
	Sepolia  Network = "sepolia"
	Arbitrum Network = "arbitrum"
	ArbTest  Network = "arb-sepolia"
)

var TestNetWork = map[Network]bool{}

func init() {
	TestNetWork[Sepolia] = true
}
