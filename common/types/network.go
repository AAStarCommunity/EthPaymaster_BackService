package types

type NetworkInfo struct {
	Name     string    `json:"main_net_name"`
	RpcUrl   string    `json:"main_net_rpc_url"`
	GasToken TokenType `json:"gas_token"`
}

type Network string

const (
	ETHEREUM_MAINNET Network = "ethereum-mainnet"
	ETHEREUM_SEPOLIA Network = "ethereum-sepolia"
	OPTIMISM_MAINNET Network = "optimism-mainnet"
	OPTIMISM_SEPOLIA Network = "optimism-sepolia"
	ARBITRUM_ONE     Network = "arbitrum-one"
	ARBITRUM_NOVA    Network = "arbitrum-nova"
	ARBITRUM_SPEOLIA Network = "arbitrum-sepolia"
	SCROLL_MAINNET   Network = "scroll-mainnet"
	SCROLL_SEPOLIA   Network = "scroll-sepolia"
	STARKET_MAINNET  Network = "starknet-mainnet"
	STARKET_SEPOLIA  Network = "starknet-sepolia"
	Base             Network = "base-mainnet"
	BaseSepolia      Network = "base-sepolia"
)

type NewWorkStack string

const (
	OPSTACK       NewWorkStack = "opstack"
	ARBSTACK      NewWorkStack = "arbstack"
	DEFAULT_STACK NewWorkStack = "default"
)
