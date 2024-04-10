package types

type NetworkInfo struct {
	Name     string    `json:"main_net_name"`
	RpcUrl   string    `json:"main_net_rpc_url"`
	GasToken TokenType `json:"gas_token"`
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
	Ethereum        Network = "ethereum"
	Sepolia         Network = "sepolia"
	Optimism        Network = "optimism"
	Optimismsepolia Network = "optimism-sepolia"
	ArbitrumOne     Network = "arbitrum-one"
	ArbitrumSeplia  Network = "arbitrum-sepolia"
	ScrollSepolia   Network = "scroll-sepolia"
	Starknet        Network = "starknet"
)
