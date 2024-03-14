package paymaster_pay_type

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
)

var (
	PaymasterDataGeneratorFactories map[types.PayType]PaymasterPayTypeExecutor
)

func init() {
	PaymasterDataGeneratorFactories = make(map[types.PayType]PaymasterPayTypeExecutor)
	PaymasterDataGeneratorFactories[types.PayTypeVerifying] = &VerifyingPaymasterExecutor{}
	PaymasterDataGeneratorFactories[types.PayTypeERC20] = &Erc20PaymasterExecutor{}
}

type PaymasterPayTypeExecutor interface {
	GeneratePayMasterAndData(strategy *model.Strategy, userOp *model.UserOperation, gasResponse *model.ComputeGasResponse, extra map[string]any) ([]byte, error)
	ValidateGas(userOp *model.UserOperation, response *model.ComputeGasResponse, strategy *model.Strategy) error
}

func GetPaymasterDataExecutor(payType types.PayType) PaymasterPayTypeExecutor {
	paymasterDataGenerator, ok := PaymasterDataGeneratorFactories[payType]
	if !ok {
		return nil
	}
	return paymasterDataGenerator
}
