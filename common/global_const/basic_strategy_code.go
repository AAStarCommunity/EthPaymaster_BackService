package global_const

type BasicStrategyCode string

const (
	StrategyCodeEthereumSepoliaV06Verify BasicStrategyCode = "Ethereum_Sepolia_v06_verifyPaymaster"
	StrategyCodeOptimismSepoliaV06Verify BasicStrategyCode = "Optimism_Sepolia_v06_verifyPaymaster"
	StrategyCodeArbitrumSpeoliaV06Verify BasicStrategyCode = "Arbitrum_Sepolia_v06_verifyPaymaster"
	StrategyCodeScrollSepoliaV06Verify   BasicStrategyCode = "Scroll_Sepolia_v06_verifyPaymaster"
	StrategyCodeBaseSepoliaV06Verify     BasicStrategyCode = "Base_Sepolia_v06_verifyPaymaster"

	StrategyCodeEthereumSepoliaV06Erc20 BasicStrategyCode = "Ethereum_Sepolia_v06_erc20Paymaster"
	StrategyCodeOptimismSepoliaV06Erc20 BasicStrategyCode = "Optimism_Sepolia_v06_erc20Paymaster"
	StrategyCodeArbitrumSpeoliaV06Erc20 BasicStrategyCode = "Arbitrum_Sepolia_v06_erc20Paymaster"
	StrategyCodeScrollSepoliaV06Erc20   BasicStrategyCode = "Scroll_Sepolia_v06_erc20Paymaster"
	StrategyCodeBaseSepoliaV06Erc20     BasicStrategyCode = "Base_Sepolia_v06_erc20Paymaster"
)
