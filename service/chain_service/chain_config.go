package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
)

var NetworkInfoMap map[network.Network]*network.NetworkInfo

func init() {
	ConfigInit()
}
func ConfigInit() {
	//TODO api key secret store
	NetworkInfoMap = map[network.Network]*network.NetworkInfo{
		network.Ethereum: {
			Name:     "ethereum",
			RpcUrl:   "https://eth-mainnet.g.alchemy.com/v2/bIZQS43-rJMgv2_SiHqfVvXa-Z1UGoGt",
			GasToken: types.ETH,
		},
		network.Sepolia: {
			Name:     "sepolia",
			RpcUrl:   "https://eth-sepolia.g.alchemy.com/v2/wKeLycGxgYRykgf0aGfcpEkUtqyLQg4v",
			GasToken: types.ETH,
		},
	}
}
