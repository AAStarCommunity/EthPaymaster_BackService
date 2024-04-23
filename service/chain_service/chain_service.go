package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"github.com/ethereum/go-ethereum/common"
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
	} else if chain == types.StarketMainnet || chain == types.StarketSepolia {
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
func GetPreVerificationGas(chain types.Network, userOp *userop.UserOpInput, strategy *model.Strategy, gasFeeResult *model.GasPrice) (*big.Int, error) {
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
	preGas = preGas.Mul(preGas, types.HUNDRED_PLUS_ONE_BIGINT)
	preGas = preGas.Div(preGas, types.HUNDRED_BIGINT)
	return preGas, nil
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
func SimulateHandleOp(networkParam types.Network, op *userop.UserOpInput, strategy *model.Strategy) (*model.SimulateHandleOpResult, error) {
	executor := network.GetEthereumExecutor(networkParam)
	entrypointVersion := strategy.GetStrategyEntryPointVersion()
	if entrypointVersion == types.EntryPointV07 {

		return executor.SimulateV06HandleOp(op, strategy.GetEntryPointAddress())

	} else if entrypointVersion == types.EntrypointV06 {
		return executor.SimulateV07HandleOp(op, strategy.GetEntryPointAddress())
	}
	return nil, xerrors.Errorf("[never be here]entrypoint version %s not support", entrypointVersion)
	//TODO Starknet
}
