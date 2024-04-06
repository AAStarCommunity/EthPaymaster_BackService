package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"math"
)

const balanceOfAbi = `[{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`

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
func GetAddressTokenBalance(networkParam network.Network, address common.Address, tokenTypeParam types.TokenType) (float64, error) {
	executor := network.GetEthereumExecutor(networkParam)
	bananceResult, err := executor.GetUserTokenBalance(address, tokenTypeParam)
	if err != nil {
		return 0, err
	}

	balanceResultFloat := float64(bananceResult.Int64()) * math.Pow(10, -6)
	return balanceResultFloat, nil

}
