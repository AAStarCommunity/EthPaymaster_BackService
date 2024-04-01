package gas_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/erc20_token"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/paymaster_pay_type"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"github.com/ethereum/go-ethereum"
	"golang.org/x/xerrors"
	"math/big"
)

func ComputeGas(userOp *userop.UserOperation, strategy *model.Strategy) (*model.ComputeGasResponse, error) {
	gasPrice, gasPriceErr := chain_service.GetGasPrice(strategy.GetNewWork())
	//TODO calculate the maximum possible fee the account needs to pay (based on validation and call gas limits, and current gas values)
	if gasPriceErr != nil {
		return nil, gasPriceErr
	}
	estimateCallGasLimit, _ := chain_service.EstimateGasLimitAndCost(strategy.GetNewWork(), ethereum.CallMsg{
		From: strategy.GetEntryPointAddress(),
		To:   &userOp.Sender,
		Data: userOp.CallData,
	})
	userOpCallGasLimit := userOp.CallGasLimit.Uint64()
	if estimateCallGasLimit > userOpCallGasLimit*12/10 {
		return nil, xerrors.Errorf("estimateCallGasLimit %d > userOpCallGasLimit %d", estimateCallGasLimit, userOpCallGasLimit)
	}

	payMasterPostGasLimit := GetPayMasterGasLimit()
	maxGasLimit := big.NewInt(0).Add(userOp.CallGasLimit, userOp.VerificationGasLimit)
	maxGasLimit = maxGasLimit.Add(maxGasLimit, payMasterPostGasLimit)
	maxFee := new(big.Int).Mul(maxGasLimit, gasPrice.MaxBasePriceWei)
	maxFeePriceInEther := new(big.Float).SetInt(maxFee)
	maxFeePriceInEther.Quo(maxFeePriceInEther, chain_service.EthWeiFactor)
	tokenCost, err := getTokenCost(strategy, maxFeePriceInEther)
	if err != nil {
		return nil, err
	}
	var usdCost float64
	if erc20_token.IsStableToken(strategy.GetUseToken()) {
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
func GetPayMasterGasLimit() *big.Int {
	//TODO
	return big.NewInt(0)
}

func ValidateGas(userOp *userop.UserOperation, gasComputeResponse *model.ComputeGasResponse, strategy *model.Strategy) error {
	paymasterDataExecutor := paymaster_pay_type.GetPaymasterDataExecutor(strategy.GetPayType())
	if paymasterDataExecutor == nil {
		return xerrors.Errorf(" %s paymasterDataExecutor not found", strategy.GetPayType())
	}
	return paymasterDataExecutor.ValidateGas(userOp, gasComputeResponse, strategy)
}
