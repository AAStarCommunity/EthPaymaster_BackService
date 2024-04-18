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
func GetPreVerificationGas(chain types.Network, userOp *userop.BaseUserOp, strategy *model.Strategy, gasFeeResult *model.GasFeePerGasResult) (*big.Int, error) {
	stack := conf.GetNetWorkStack(chain)
	preGasFunc := network.PreVerificationGasFuncMap[stack]
	return preGasFunc(userOp, strategy, gasFeeResult)
}

func GetEntryPointDeposit(entrypoint string, depositAddress string) uint256.Int {

	return uint256.Int{1}
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
func SimulateHandleOp(networkParam types.Network, op *userop.BaseUserOp, strategy model.Strategy) (*model.SimulateHandleOpResult, error) {
	executor := network.GetEthereumExecutor(networkParam)
	opValue := *op
	entrypointVersion := opValue.GetEntrypointVersion()
	if entrypointVersion == types.EntryPointV07 {
		userOpV6 := opValue.(*userop.UserOperationV06)
		return executor.SimulateV06HandleOp(userOpV6, strategy.GetEntryPointAddress())

	} else if entrypointVersion == types.EntrypointV06 {
		userOpV7 := opValue.(*userop.UserOperationV07)
		return executor.SimulateV07HandleOp(userOpV7, strategy.GetEntryPointAddress())
	}
	return nil, xerrors.Errorf("[never be here]entrypoint version %s not support", entrypointVersion)
	//TODO Starknet
}
func GetVertificationGasLimit(chain types.Network) (*big.Int, error) {
	return nil, nil
}
