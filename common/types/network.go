package types

type NetworkInfo struct {
	Name     string    `json:"main_net_name"`
	RpcUrl   string    `json:"main_net_rpc_url"`
	GasToken TokenType `json:"gas_token"`
}

type Network string

const (
	Ethereum        Network = "ethereum"
	Sepolia         Network = "sepolia"
	Optimism        Network = "optimism"
	OptimismSepolia Network = "optimism-sepolia"
	ArbitrumOne     Network = "arbitrum-one"
	ArbitrumSeplia  Network = "arbitrum-sepolia"
	Scroll          Network = "scroll"
	ScrollSepolia   Network = "scroll-sepolia"
	Starknet        Network = "starknet"
	StarknetSepolia Network = "starknet-sepolia"
	Base            Network = "base"
	BaseSepolia     Network = "base-sepolia"
)
