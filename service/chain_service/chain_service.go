package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/xerrors"
	"math"
	"math/big"
)

func CheckContractAddressAccess(contract *common.Address, chain global_const.Network) (bool, error) {
	//todo needcache
	executor := network.GetEthereumExecutor(chain)
	return executor.CheckContractAddressAccess(contract)
}

// GetGasPrice return gas price in wei, gwei, ether
func GetGasPrice(chain global_const.Network) (*model.GasPrice, error) {
	if conf.IsEthereumAdaptableNetWork(chain) {
		ethereumExecutor := network.GetEthereumExecutor(chain)
		return ethereumExecutor.GetGasPrice()
	} else if chain == global_const.StarketMainnet || chain == global_const.StarketSepolia {
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

// GetPreVerificationGas https://github.com/eth-infinitism/bundler/blob/main/packages/sdk/src/calcPreVerificationGas.ts
func GetPreVerificationGas(userOp *user_op.UserOpInput, strategy *model.Strategy, gasFeeResult *model.GasPrice) (*big.Int, error) {
	chain := strategy.GetNewWork()
	stack := conf.GetNetWorkStack(chain)
	preGasFunc, err := network.GetPreVerificationGasFunc(stack)
	if err != nil {
		return nil, err
	}
	preGas, err := preGasFunc(userOp, strategy, gasFeeResult)
	if err != nil {
		return nil, err
	}
	// add 10% buffer
	preGas = preGas.Mul(preGas, global_const.HundredPlusOneBigint)
	preGas = preGas.Div(preGas, global_const.HundredBigint)
	return preGas, nil
}

func GetAddressTokenBalance(networkParam global_const.Network, address common.Address, tokenTypeParam global_const.TokenType) (float64, error) {
	executor := network.GetEthereumExecutor(networkParam)
	bananceResult, err := executor.GetUserTokenBalance(address, tokenTypeParam)
	if err != nil {
		return 0, err
	}

	balanceResultFloat := float64(bananceResult.Int64()) * math.Pow(10, -6)
	return balanceResultFloat, nil

}
func SimulateHandleOp(op *user_op.UserOpInput, strategy *model.Strategy) (*model.SimulateHandleOpResult, error) {
	networkParam := strategy.GetNewWork()
	executor := network.GetEthereumExecutor(networkParam)
	entrypointVersion := strategy.GetStrategyEntrypointVersion()
	if entrypointVersion == global_const.EntrypointV06 {

		return executor.SimulateV06HandleOp(*op, strategy.GetEntryPointAddress())

	} else if entrypointVersion == global_const.EntryPointV07 {
		return executor.SimulateV07HandleOp(*op, strategy.GetEntryPointAddress())
	}
	return nil, xerrors.Errorf("[never be here]entrypoint version %s not support", entrypointVersion)
	//TODO Starknet
}
