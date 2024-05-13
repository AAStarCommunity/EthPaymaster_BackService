package global_const

type NetworkInfo struct {
	Name     string    `json:"main_net_name"`
	RpcUrl   string    `json:"main_net_rpc_url"`
	GasToken TokenType `json:"gas_token"`
}

type Network string

const (
	EthereumMainnet Network = "ethereum-mainnet"
	EthereumSepolia Network = "ethereum-sepolia"
	OptimismMainnet Network = "optimism-mainnet"
	OptimismSepolia Network = "optimism-sepolia"
	ArbitrumOne     Network = "arbitrum-one"
	ArbitrumNova    Network = "arbitrum-nova"
	ArbitrumSpeolia Network = "arbitrum-sepolia"
	ScrollMainnet   Network = "scroll-mainnet"
	ScrollSepolia   Network = "scroll-sepolia"
	StarketMainnet  Network = "starknet-mainnet"
	StarketSepolia  Network = "starknet-sepolia"
	BaseMainnet     Network = "base-mainnet"
	BaseSepolia     Network = "base-sepolia"
)

type NewWorkStack string

const (
	OpStack      NewWorkStack = "opstack"
	ArbStack     NewWorkStack = "arbstack"
	DefaultStack NewWorkStack = "default"
)
