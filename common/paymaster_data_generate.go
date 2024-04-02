package common

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
)

var GenerateFuncMap = map[types.PayType]GeneratePaymasterDataFunc{}

func init() {
	GenerateFuncMap[types.PayTypeVerifying] = GenerateBasicPaymasterData()
	GenerateFuncMap[types.PayTypeERC20] = GenerateBasicPaymasterData()
	GenerateFuncMap[types.PayTypeSuperVerifying] = GenerateSuperContractPaymasterData()
}

type GeneratePaymasterDataFunc = func(strategy *model.Strategy, userOp *userop.BaseUserOp, gasResponse *model.ComputeGasResponse) (string, string, error)

func GenerateBasicPaymasterData() GeneratePaymasterDataFunc {
	return func(strategy *model.Strategy, userOp *userop.BaseUserOp, gasResponse *model.ComputeGasResponse) (string, string, error) {
		return "", "", nil
	}
}

func GenerateSuperContractPaymasterData() GeneratePaymasterDataFunc {
	return func(strategy *model.Strategy, userOp *userop.BaseUserOp, gasResponse *model.ComputeGasResponse) (string, string, error) {
		return "", "", nil
	}
}
