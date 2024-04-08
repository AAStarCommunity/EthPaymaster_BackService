package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
)

var NetworkInfoMap map[types.Network]*types.NetworkInfo

func init() {
	ConfigInit()
}
func ConfigInit() {
	//TODO api key secret store
	NetworkInfoMap = map[types.Network]*types.NetworkInfo{
		types.Ethereum: {
			Name:     "ethereum",
			RpcUrl:   "https://eth-mainnet.g.alchemy.com/v2/bIZQS43-rJMgv2_SiHqfVvXa-Z1UGoGt",
			GasToken: types.ETH,
		},
		types.Sepolia: {
			Name:     "sepolia",
			RpcUrl:   "https://eth-sepolia.g.alchemy.com/v2/wKeLycGxgYRykgf0aGfcpEkUtqyLQg4v",
			GasToken: types.ETH,
		},
	}
}
