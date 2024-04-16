package types

type NetworkInfo struct {
	Name     string    `json:"main_net_name"`
	RpcUrl   string    `json:"main_net_rpc_url"`
	GasToken TokenType `json:"gas_token"`
}

type Network string

const (
	ETHEREUM_MAINNET Network = "ethereum"
	ETHEREUM_SEPOLIA Network = "sepolia"
	OPTIMISM_MAINNET Network = "optimism"
	OPTIMISM_SEPOLIA Network = "optimism-sepolia"
	ARBITRUM_ONE     Network = "arbitrum-one"
	ARBITRUM_SPEOLIA Network = "arbitrum-sepolia"
	SCROLL_MAINNET   Network = "scroll"
	SCROLL_SEPOLIA   Network = "scroll-sepolia"
	STARKET_MAINNET  Network = "starknet"
	STARKET_SEPOLIA  Network = "starknet-sepolia"
	Base             Network = "base"
	BaseSepolia      Network = "base-sepolia"
)
