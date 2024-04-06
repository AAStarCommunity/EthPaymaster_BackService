package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/tokens"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"math"
)

const balanceOfAbi = `[{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`

var TokenAddressMap map[network.Network]*map[tokens.TokenType]common.Address

func init() {
	TokenAddressMap = map[network.Network]*map[tokens.TokenType]common.Address{
		network.Ethereum: {
			tokens.ETH: common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7"),
		},
		network.Sepolia: {
			tokens.USDT: common.HexToAddress("0xaa8e23fb1079ea71e0a56f48a2aa51851d8433d0"),
			tokens.USDC: common.HexToAddress("0x1c7d4b196cb0c7b01d743fbc6116a902379c7238"),
		},
	}
}
func CheckContractAddressAccess(contract *common.Address, chain network.Network) (bool, error) {
	executor := network.GetEthereumExecutor(chain)
	return executor.CheckContractAddressAccess(contract)
}

// GetGasPrice return gas price in wei, gwei, ether
func GetGasPrice(chain network.Network) (*model.GasPrice, error) {
	ethereumExecutor := network.GetEthereumExecutor(chain)
	return ethereumExecutor.GetCurGasPrice()
	//TODO starknet
}

func GetEntryPointDeposit(entrypoint string, depositAddress string) uint256.Int {
	return uint256.Int{1}
}
func EstimateUserOpGas(strategy *model.Strategy, op *userop.BaseUserOp) (uint64, error) {
	ethereumExecutor := network.GetEthereumExecutor(strategy.GetNewWork())
	return ethereumExecutor.EstimateUserOpGas(strategy.GetEntryPointAddress(), op)
}
func GetAddressTokenBalance(networkParam network.Network, address common.Address, tokenTypeParam tokens.TokenType) (float64, error) {
	executor := network.GetEthereumExecutor(networkParam)
	bananceResult, err := executor.GetUserTokenBalance(address, tokenTypeParam)
	if err != nil {
		return 0, err
	}

	balanceResultFloat := float64(bananceResult.Int64()) * math.Pow(10, -6)
	return balanceResultFloat, nil

}
