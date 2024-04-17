package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"golang.org/x/xerrors"
	"math"
	"math/big"
)

const balanceOfAbi = `[{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`

func CheckContractAddressAccess(contract *common.Address, chain types.Network) (bool, error) {
	//todo needcache
	executor := network.GetEthereumExecutor(chain)
	return executor.CheckContractAddressAccess(contract)
}

// GetGasPrice return gas price in wei, gwei, ether
func GetGasPrice(chain types.Network) (*model.GasPrice, error) {
	if conf.IsEthereumAdaptableNetWork(chain) {
		ethereumExecutor := network.GetEthereumExecutor(chain)
		return ethereumExecutor.GetGasPrice()
	} else if chain == types.STARKET_MAINNET || chain == types.STARKET_SEPOLIA {
		starknetExecutor := network.GetStarknetExecutor()
		return starknetExecutor.GetGasPrice()
	} else {
		return nil, xerrors.Errorf("chain %s not support", chain)
	}
	//MaxFeePerGas
	//MaxPriorityPrice
	//preOpGas (get verificationGasLimit from preOpGas)
	//

}
func GetCallGasLimit(chain types.Network) (*big.Int, *big.Int, error) {
	//TODO
	return nil, nil, nil
}

// GetPreVerificationGas https://github.com/eth-infinitism/bundler/blob/main/packages/sdk/src/calcPreVerificationGas.ts
func GetPreVerificationGas(chain types.Network) (*big.Int, error) {
	stack := conf.GetNetWorkStack(chain)
	preGasFunc := network.PreVerificationGasFuncMap[stack]
	return preGasFunc()
}

func GetEntryPointDeposit(entrypoint string, depositAddress string) uint256.Int {
	return uint256.Int{1}
}
func EstimateUserOpGas(strategy *model.Strategy, op *userop.BaseUserOp) (uint64, error) {
	ethereumExecutor := network.GetEthereumExecutor(strategy.GetNewWork())
	return ethereumExecutor.EstimateUserOpCallGas(strategy.GetEntryPointAddress(), op)
}
func GetAddressTokenBalance(networkParam types.Network, address common.Address, tokenTypeParam types.TokenType) (float64, error) {
	executor := network.GetEthereumExecutor(networkParam)
	bananceResult, err := executor.GetUserTokenBalance(address, tokenTypeParam)
	if err != nil {
		return 0, err
	}

	balanceResultFloat := float64(bananceResult.Int64()) * math.Pow(10, -6)
	return balanceResultFloat, nil

}
func SimulateHandleOp(networkParam types.Network) (*model.SimulateHandleOpResult, error) {

	return nil, nil

}
func GetVertificationGasLimit(chain types.Network) (*big.Int, error) {
	return nil, nil
}
