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

type Network string

const (
	Ethereum Network = "ethereum"
	Sepolia  Network = "sepolia"
	Arbitrum Network = "arbitrum"
)

var TestNetWork = map[Network]bool{}

func init() {
	TestNetWork[Sepolia] = true
}
