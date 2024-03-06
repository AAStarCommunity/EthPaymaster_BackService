package types

type NetworkInfo struct {
	Name   string `json:"main_net_name"`
	RpcUrl string `json:"main_net_rpc_url"`
}

//type Chain string
//
//const (
//	Ethereum Chain = "Ethereum"
//	Arbitrum Chain = "Arbitrum"
//	Optimism Chain = "Optimism"
//)

type NetWork string

const (
	Ethereum NetWork = "ethereum"
	Sepolia  NetWork = "sepolia"
	Arbitrum NetWork = "arbitrum"
)
