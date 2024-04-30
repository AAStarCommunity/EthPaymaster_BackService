package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"math"
	"math/big"
)

func CheckContractAddressAccess(contract *common.Address, chain global_const.Network) (bool, error) {
	//todo needcache
	executor := network.GetEthereumExecutor(chain)
	return executor.CheckContractAddressAccess(contract)
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
func GetPaymasterEntryPointBalance(strategy *model.Strategy) (*big.Float, error) {
	networkParam := strategy.GetNewWork()
	paymasterAddress := strategy.GetPaymasterAddress()
	logrus.Debug("paymasterAddress", paymasterAddress)
	executor := network.GetEthereumExecutor(networkParam)
	balance, err := executor.GetPaymasterDeposit(paymasterAddress)
	if err != nil {
		return nil, err
	}
	logrus.Debug("balance", balance)
	balanceResultFloat := utils.ConvertBalanceToEther(balance)

	return balanceResultFloat, nil
}
func SimulateHandleOp(op *user_op.UserOpInput, strategy *model.Strategy) (*model.SimulateHandleOpResult, error) {
	networkParam := strategy.GetNewWork()
	executor := network.GetEthereumExecutor(networkParam)
	entrypointVersion := strategy.GetStrategyEntrypointVersion()
	if entrypointVersion == global_const.EntrypointV06 {

		return executor.SimulateV06HandleOp(*op, strategy.GetEntryPointAddress())

	} else if entrypointVersion == global_const.EntrypointV07 {
		return executor.SimulateV07HandleOp(*op, strategy.GetEntryPointAddress())
	}
	return nil, xerrors.Errorf("[never be here]entrypoint version %s not support", entrypointVersion)
	//TODO Starknet
}
