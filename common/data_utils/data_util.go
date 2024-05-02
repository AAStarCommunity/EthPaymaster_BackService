package data_utils

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/paymaster_data"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
)

func GetUserOpWithPaymasterAndDataForSimulate(op user_op.UserOpInput, strategy *model.Strategy, paymasterDataInput *paymaster_data.PaymasterDataInput, gasPriceResult *model.GasPrice) (*user_op.UserOpInput, error) {
	executor := network.GetEthereumExecutor(strategy.GetNewWork())
	op.MaxFeePerGas = gasPriceResult.MaxFeePerGas
	op.MaxPriorityFeePerGas = gasPriceResult.MaxPriorityFeePerGas
	paymasterDataInput.PaymasterPostOpGasLimit = global_const.DummyPaymasterPostopGaslimitBigint
	paymasterDataInput.PaymasterVerificationGasLimit = global_const.DummyPaymasterVerificationgaslimitBigint
	op.AccountGasLimits = user_op.DummyAccountGasLimits
	op.GasFees = user_op.DummyGasFees
	paymasterData, err := executor.GetPaymasterData(&op, strategy, paymasterDataInput)
	if err != nil {
		return nil, err
	}
	op.PaymasterAndData = paymasterData
	return &op, nil
}
