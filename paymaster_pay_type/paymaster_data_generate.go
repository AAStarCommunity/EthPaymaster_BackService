package paymaster_pay_type

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"fmt"
	"strconv"
)

var GenerateFuncMap = map[types.PayType]GeneratePaymasterDataFunc{}

func init() {
	GenerateFuncMap[types.PayTypeVerifying] = GenerateBasicPaymasterData()
	GenerateFuncMap[types.PayTypeERC20] = GenerateBasicPaymasterData()
	GenerateFuncMap[types.PayTypeSuperVerifying] = GenerateSuperContractPaymasterData()
}

type GeneratePaymasterDataFunc = func(strategy *model.Strategy, userOp *user_op.UserOpInput, gasResponse *model.ComputeGasResponse) (string, error)

func GenerateBasicPaymasterData() GeneratePaymasterDataFunc {
	return func(strategy *model.Strategy, userOp *user_op.UserOpInput, gasResponse *model.ComputeGasResponse) (string, error) {
		validStart, validEnd := getValidTime(strategy)
		message := fmt.Sprintf("%s%s%s", strategy.GetPaymasterAddress().String(), validEnd, validStart)
		return message, nil
	}
}

func GenerateSuperContractPaymasterData() GeneratePaymasterDataFunc {
	return func(strategy *model.Strategy, userOp *user_op.UserOpInput, gasResponse *model.ComputeGasResponse) (string, error) {
		validStart, validEnd := getValidTime(strategy)
		message := fmt.Sprintf("%s%s%s", strategy.GetPaymasterAddress().String(), validEnd, validStart)
		return message, nil
	}
}

func getValidTime(strategy *model.Strategy) (string, string) {

	currentTimestampStr := strconv.FormatInt(strategy.ExecuteRestriction.EffectiveStartTime.Int64(), 16)
	futureTimestampStr := strconv.FormatInt(strategy.ExecuteRestriction.EffectiveEndTime.Int64(), 16)
	currentTimestampStrSupply := utils.SupplyZero(currentTimestampStr, 64)
	futureTimestampStrSupply := utils.SupplyZero(futureTimestampStr, 64)
	return currentTimestampStrSupply, futureTimestampStrSupply
}
