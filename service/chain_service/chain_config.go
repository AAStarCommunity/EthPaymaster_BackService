package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/token"
	"github.com/ethereum/go-ethereum/ethclient"
)

var NetworkInfoMap map[network.Network]*network.NetworkInfo
var EthCompatibleNetWorkClientMap map[network.Network]*ethclient.Client

func init() {
	ConfigInit()
	ClientInit()
}
func ConfigInit() {
	//TODO api key secret store
	NetworkInfoMap = map[network.Network]*network.NetworkInfo{
		network.Ethereum: {
			Name:     "ethereum",
			RpcUrl:   "https://eth-mainnet.g.alchemy.com/v2/bIZQS43-rJMgv2_SiHqfVvXa-Z1UGoGt",
			GasToken: token.ETH,
		},
		network.Sepolia: {
			Name:     "sepolia",
			RpcUrl:   "https://eth-sepolia.g.alchemy.com/v2/wKeLycGxgYRykgf0aGfcpEkUtqyLQg4v",
			GasToken: token.ETH,
		},
	}
}

func ClientInit() {
	EthCompatibleNetWorkClientMap = make(map[network.Network]*ethclient.Client)
	for chain, networkInfo := range NetworkInfoMap {
		client, err := ethclient.Dial(networkInfo.RpcUrl)
		if err != nil {
			panic(err)
		}
		EthCompatibleNetWorkClientMap[chain] = client
		continue
	}
}
