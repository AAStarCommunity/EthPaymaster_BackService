package data_utils

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/paymaster_data"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"github.com/sirupsen/logrus"
)

func GetUserOpWithPaymasterAndDataForSimulate(op user_op.UserOpInput, strategy *model.Strategy, paymasterDataInput *paymaster_data.PaymasterDataInput, gasPriceResult *model.GasPrice) (*user_op.UserOpInput, error) {
	executor := network.GetEthereumExecutor(strategy.GetNewWork())
	op.MaxFeePerGas = gasPriceResult.MaxFeePerGas
	op.MaxPriorityFeePerGas = gasPriceResult.MaxPriorityFeePerGas
	logrus.Debug("MaxFeePerGas", op.MaxFeePerGas)
	logrus.Debug("MaxPriorityFeePerGas", op.MaxPriorityFeePerGas)
	paymasterDataInput.PaymasterPostOpGasLimit = global_const.DummyPaymasterPostoperativelyBigint
	paymasterDataInput.PaymasterVerificationGasLimit = global_const.DummyPaymasterOversimplificationBigint
	op.AccountGasLimits = user_op.DummyAccountGasLimits
	op.GasFees = user_op.DummyGasFees
	if op.PreVerificationGas == nil {
		op.PreVerificationGas = global_const.DummyReverificationsBigint
	}
	if op.VerificationGasLimit == nil {
		op.VerificationGasLimit = global_const.DummyVerificationGasLimit
	}
	if op.Signature == nil {
		op.Signature = global_const.DummySignatureByte
	}
	if op.CallGasLimit == nil {
		op.CallGasLimit = global_const.DummyCallGasLimit
	}

	paymasterData, _, err := executor.GetPaymasterData(&op, strategy, paymasterDataInput)
	if err != nil {
		return nil, err
	}
	op.PaymasterAndData = paymasterData
	return &op, nil
}
