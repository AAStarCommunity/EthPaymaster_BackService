package gas_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/paymaster_pay_type"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"golang.org/x/xerrors"
	"math/big"
)

func ComputeGas(userOp *userop.BaseUserOp, strategy *model.Strategy) (*model.ComputeGasResponse, error) {
	gasPrice, gasPriceErr := chain_service.GetGasPrice(strategy.GetNewWork())
	//TODO calculate the maximum possible fee the account needs to pay (based on validation and call gas limits, and current gas values)
	if gasPriceErr != nil {
		return nil, gasPriceErr
	}
	userOpValue := *userOp
	var maxFeePriceInEther *big.Float
	var maxFee *big.Int
	estimateCallGasLimit, err := chain_service.EstimateUserOpGas(strategy, userOp)
	if err != nil {
		return nil, err
	}
	switch userOpValue.GetEntrypointVersion() {
	case types.EntrypointV06:
		{
			useropV6Value := userOpValue.(*userop.UserOperationV06)
			userOpCallGasLimit := useropV6Value.CallGasLimit.Uint64()
			if estimateCallGasLimit > userOpCallGasLimit*12/10 {
				return nil, xerrors.Errorf("estimateCallGasLimit %d > userOpCallGasLimit %d", estimateCallGasLimit, userOpCallGasLimit)
			}

			payMasterPostGasLimit := GetPayMasterGasLimit()
			maxGasLimit := big.NewInt(0).Add(useropV6Value.CallGasLimit, useropV6Value.VerificationGasLimit)
			maxGasLimit = maxGasLimit.Add(maxGasLimit, payMasterPostGasLimit)
			maxFee = new(big.Int).Mul(maxGasLimit, gasPrice.MaxBasePriceWei)
			maxFeePriceInEther = new(big.Float).SetInt(maxFee)
			maxFeePriceInEther.Quo(maxFeePriceInEther, network.EthWeiFactor)
		}
		break
	case types.EntryPointV07:
		{

		}
		break

	}

	tokenCost, err := getTokenCost(strategy, maxFeePriceInEther)
	if err != nil {
		return nil, err
	}
	var usdCost float64
	if types.IsStableToken(strategy.GetUseToken()) {
		usdCost, _ = tokenCost.Float64()
	} else {
		usdCost, _ = utils.GetPriceUsd(strategy.GetUseToken())
	}

	// TODO get PaymasterCallGasLimit
	return &model.ComputeGasResponse{
		GasInfo:   gasPrice,
		TokenCost: tokenCost,
		Network:   strategy.GetNewWork(),
		Token:     strategy.GetUseToken(),
		UsdCost:   usdCost,
		MaxFee:    *maxFee,
	}, nil
}

func getTokenCost(strategy *model.Strategy, tokenCount *big.Float) (*big.Float, error) {
	if strategy.GetPayType() == types.PayTypeERC20 {
		formTokenType := chain_service.NetworkInfoMap[strategy.GetNewWork()].GasToken
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
func GetPayMasterGasLimit() *big.Int {
	//TODO
	return big.NewInt(0)
}

func ValidateGas(userOp *userop.BaseUserOp, gasComputeResponse *model.ComputeGasResponse, strategy *model.Strategy) error {
	validateFunc := paymaster_pay_type.GasValidateFuncMap[strategy.GetPayType()]
	err := validateFunc(userOp, gasComputeResponse, strategy)
	if err != nil {
		return err
	}
	return nil
}
