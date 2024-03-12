package paymaster_data_generator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
)

var (
	PaymasterDataGeneratorFactories map[types.PayType]PaymasterDataGenerator
)

func init() {
	PaymasterDataGeneratorFactories = make(map[types.PayType]PaymasterDataGenerator)
	PaymasterDataGeneratorFactories[types.PayTypeVerifying] = &VerifyingPaymasterGenerator{}
	PaymasterDataGeneratorFactories[types.PayTypeERC20] = &Erc20PaymasterGenerator{}
}

type PaymasterDataGenerator interface {
	GeneratePayMaster(strategy *model.Strategy, userOp *model.UserOperation, gasResponse *model.ComputeGasResponse, extra map[string]any) ([]byte, error)
}

func GetPaymasterDataGenerator(payType types.PayType) PaymasterDataGenerator {

	paymasterDataGenerator, ok := PaymasterDataGeneratorFactories[payType]
	if !ok {
		return nil
	}
	return paymasterDataGenerator
}
