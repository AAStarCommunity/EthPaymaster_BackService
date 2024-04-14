package gas_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"AAStarCommunity/EthPaymaster_BackService/paymaster_pay_type"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"golang.org/x/xerrors"
	"math/big"
)

// https://blog.particle.network/bundler-predicting-gas/
func ComputeGas(userOp userop.BaseUserOp, strategy *model.Strategy) (*model.ComputeGasResponse, *userop.BaseUserOp, error) {
	gasPrice, gasPriceErr := chain_service.GetGasPrice(strategy.GetNewWork())
	//TODO calculate the maximum possible fee the account needs to pay (based on validation and call gas limits, and current gas values)
	if gasPriceErr != nil {
		return nil, nil, gasPriceErr
	}
	paymasterUserOp := userOp
	var maxFeePriceInEther *big.Float
	var maxFee *big.Int
	opEstimateGas := model.UserOpEstimateGas{}
	switch paymasterUserOp.GetEntrypointVersion() {
	case types.EntrypointV06:
		{
			// Get MaxFeePerGas And MaxPriorityFeePerGas
			// if MaxFeePerGas <=0 use recommend gas price
			useropV6Value := paymasterUserOp.(*userop.UserOperationV06)
			opEstimateGas.MaxFeePerGas = useropV6Value.MaxFeePerGas
			opEstimateGas.MaxPriorityFeePerGas = useropV6Value.MaxPriorityFeePerGas
			if utils.IsLessThanZero(useropV6Value.MaxFeePerGas) {
				opEstimateGas.MaxFeePerGas = gasPrice.MaxBasePriceWei
			}
			if utils.IsLessThanZero(useropV6Value.MaxPriorityFeePerGas) {
				opEstimateGas.MaxPriorityFeePerGas = gasPrice.MaxPriorityPriceWei
			}
			// TODO Get verificationGasLimit callGasLimit
			//estimateCallGasLimit, err := chain_service.EstimateUserOpCallGas(strategy, userOp)
			//if err != nil {
			//	return nil, err
			//}
			//if estimateCallGasLimit > userOpCallGasLimit*12/10 {
			//	return nil, xerrors.Errorf("estimateCallGasLimit %d > userOpCallGasLimit %d", estimateCallGasLimit, userOpCallGasLimit)
			//}
			useropV6Value.MaxFeePerGas = opEstimateGas.MaxFeePerGas
			useropV6Value.MaxPriorityFeePerGas = opEstimateGas.MaxPriorityFeePerGas

			verficationGasLimit, callGasLimit, err := useropV6Value.EstimateGasLimit(strategy)
			if err != nil {
				return nil, nil, err
			}
			opEstimateGas.VerificationGasLimit = big.NewInt(int64(verficationGasLimit))
			opEstimateGas.CallGasLimit = big.NewInt(int64(callGasLimit))

			// TODO  Get PreVerificationGas

			// over UserOp

			useropV6Value.VerificationGasLimit = opEstimateGas.VerificationGasLimit
			useropV6Value.CallGasLimit = opEstimateGas.CallGasLimit
		}
		break
	case types.EntryPointV07:
		{
			useropV7Value := paymasterUserOp.(*userop.UserOperationV07)
			useropV7Value.PaymasterVerificationGasLimit = opEstimateGas.PaymasterVerificationGasLimit
		}
		break

	}

	tokenCost, err := getTokenCost(strategy, maxFeePriceInEther)
	if err != nil {
		return nil, nil, err
	}
	var usdCost float64
	if types.IsStableToken(strategy.GetUseToken()) {
		usdCost, _ = tokenCost.Float64()
	} else {
		usdCost, _ = utils.GetPriceUsd(strategy.GetUseToken())
	}

	// TODO get PaymasterCallGasLimit
	return &model.ComputeGasResponse{
		GasInfo:       gasPrice,
		TokenCost:     tokenCost,
		OpEstimateGas: &opEstimateGas,
		Network:       strategy.GetNewWork(),
		Token:         strategy.GetUseToken(),
		UsdCost:       usdCost,
		MaxFee:        *maxFee,
	}, &paymasterUserOp, nil
}

func getTokenCost(strategy *model.Strategy, tokenCount *big.Float) (*big.Float, error) {
	if strategy.GetPayType() == types.PayTypeERC20 {

		formTokenType := conf.GetGasToken(strategy.GetNewWork())
		toTokenType := strategy.GetUseToken()
		toTokenPrice, err := utils.GetToken(formTokenType, toTokenType)
		if err != nil {
			return nil, err
		}
		if toTokenPrice == 0 {
			return nil, xerrors.Errorf("toTokenPrice can not be 0")
		}
		tokenCost := new(big.Float).Mul(tokenCount, big.NewFloat(toTokenPrice))
		return tokenCost, nil
	}
	return tokenCount, nil

}

func ValidateGas(userOp *userop.BaseUserOp, gasComputeResponse *model.ComputeGasResponse, strategy *model.Strategy) error {
	validateFunc := paymaster_pay_type.GasValidateFuncMap[strategy.GetPayType()]
	err := validateFunc(userOp, gasComputeResponse, strategy)
	if err != nil {
		return err
	}
	return nil
}
