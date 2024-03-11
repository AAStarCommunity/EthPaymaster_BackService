package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var NetworkInfoMap map[types.Network]*types.NetworkInfo
var NetWorkClientMap map[types.Network]*ethclient.Client

func init() {
	ConfigInit()
	ClientInit()
}
func ConfigInit() {
	//TODO api key secret store
	NetworkInfoMap = map[types.Network]*types.NetworkInfo{
		types.Ethereum: {
			Name:   "ethereum",
			RpcUrl: "https://eth-mainnet.g.alchemy.com/v2/bIZQS43-rJMgv2_SiHqfVvXa-Z1UGoGt",
		},
		types.Sepolia: {
			Name:   "sepolia",
			RpcUrl: "https://eth-sepolia.g.alchemy.com/v2/wKeLycGxgYRykgf0aGfcpEkUtqyLQg4v",
		},
	}
}

func ClientInit() {
	NetWorkClientMap = make(map[types.Network]*ethclient.Client)
	for chain, networkInfo := range NetworkInfoMap {
		client, err := ethclient.Dial(networkInfo.RpcUrl)
		if err != nil {
			panic(err)
		}
		NetWorkClientMap[chain] = client
		continue
	}
}
