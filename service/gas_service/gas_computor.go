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
func ComputeGas(userOp *userop.BaseUserOp, strategy *model.Strategy) (*model.ComputeGasResponse, *userop.BaseUserOp, error) {
	gasPrice, gasPriceErr := chain_service.GetGasPrice(strategy.GetNewWork())
	//TODO calculate the maximum possible fee the account needs to pay (based on validation and call gas limits, and current gas values)
	if gasPriceErr != nil {
		return nil, nil, gasPriceErr
	}
	paymasterUserOp := *userOp
	var maxFeePriceInEther *big.Float
	var maxFee *big.Int
	opEstimateGas := model.UserOpEstimateGas{}

	entryPointVersion := paymasterUserOp.GetEntrypointVersion()

	verficationGasLimit, err := EstimateVerificationGasLimit(strategy)
	callGasLimit, err := EstimateCallGasLimit(strategy)
	maxFeePerGas, maxPriorityFeePerGas, baseFee := GetFeePerGas(strategy)
	opEstimateGas.VerificationGasLimit = verficationGasLimit
	opEstimateGas.CallGasLimit = callGasLimit
	opEstimateGas.MaxFeePerGas = maxFeePerGas
	opEstimateGas.MaxPriorityFeePerGas = maxPriorityFeePerGas
	opEstimateGas.BaseFee = baseFee
	preVerificationGas, err := chain_service.GetPreVerificationGas(strategy.GetNewWork(), userOp, strategy, opEstimateGas)
	opEstimateGas.PreVerificationGas = preVerificationGas
	if entryPointVersion == types.EntryPointV07 {
		opEstimateGas.PaymasterPostOpGasLimit = types.DUMMY_PAYMASTER_POSTOP_GASLIMIT_BIGINT
		opEstimateGas.PaymasterVerificationGasLimit = types.DUMMY_PAYMASTER_VERIFICATIONGASLIMIT_BIGINT
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

func GetFeePerGas(strategy *model.Strategy) (*big.Int, *big.Int, *big.Int) {
	return nil, nil, nil
}

func EstimateCallGasLimit(strategy *model.Strategy) (*big.Int, error) {
	//TODO
	return nil, nil
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

func EstimateVerificationGasLimit(strategy *model.Strategy) (*big.Int, error) {
	//TODO
	return nil, nil
}
