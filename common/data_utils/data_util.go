package data_utils

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/paymaster_data"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
)

func GetUserOpWithPaymasterAndDataForSimulate(op user_op.UserOpInput, strategy *model.Strategy, paymasterDataInput *paymaster_data.PaymasterData, gasPriceResult *model.GasPrice) (*user_op.UserOpInput, error) {
	executor := network.GetEthereumExecutor(strategy.GetNewWork())
	op.MaxFeePerGas = gasPriceResult.MaxFeePerGas
	op.MaxPriorityFeePerGas = gasPriceResult.MaxPriorityFeePerGas
	paymasterData, err := executor.GetPaymasterData(&op, strategy, paymasterDataInput)
	if err != nil {
		return nil, err
	}
	op.PaymasterAndData = paymasterData
	return &op, nil
}
