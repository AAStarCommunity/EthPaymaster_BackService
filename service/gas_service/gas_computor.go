package gas_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/paymaster_pay_type"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/xerrors"
	"math/big"
)

func ComputeGas(userOp *model.UserOperation, strategy *model.Strategy) (*model.ComputeGasResponse, error) {
	gasPrice, gasPriceErr := chain_service.GetGasPrice(strategy.NetWork)
	//TODO calculate the maximum possible fee the account needs to pay (based on validation and call gas limits, and current gas values)
	if gasPriceErr != nil {
		return nil, gasPriceErr
	}
	estimateCallGasLimit, _ := chain_service.EstimateGasLimitAndCost(strategy.NetWork, ethereum.CallMsg{
		From: common.HexToAddress(strategy.EntryPointAddress),
		To:   &userOp.Sender,
		Data: userOp.CallData,
	})
	userOpCallGasLimit := userOp.CallGasLimit.Uint64()
	if estimateCallGasLimit > userOpCallGasLimit {
		return nil, xerrors.Errorf("estimateCallGasLimit %d > userOpCallGasLimit %d", estimateCallGasLimit, userOpCallGasLimit)
	}

	payMasterPostGasLimit := GetPayMasterGasLimit()
	maxGasLimit := big.NewInt(0).Add(userOp.CallGasLimit, userOp.VerificationGasLimit)
	maxGasLimit = maxGasLimit.Add(maxGasLimit, payMasterPostGasLimit)
	maxFee := new(big.Int).Mul(maxGasLimit, gasPrice.MaxBasePriceWei)
	maxFeePriceInEther := new(big.Float).SetInt(maxFee)
	maxFeePriceInEther.Quo(maxFeePriceInEther, chain_service.EthWeiFactor)
	tokenCost, _ := getTokenCost(strategy, maxFeePriceInEther)
	var usdCost float64
	if types.IsStableToken(strategy.Token) {
		usdCost, _ = tokenCost.Float64()
	} else {
		usdCost, _ = utils.GetPriceUsd(strategy.Token)
	}

	// TODO get PaymasterCallGasLimit
	return &model.ComputeGasResponse{
		GasInfo:    gasPrice,
		TokenCost:  tokenCost,
		Network:    strategy.NetWork,
		Token:      strategy.Token,
		UsdCost:    usdCost,
		BlobEnable: strategy.Enable4844,
		MaxFee:     *maxFee,
	}, nil
}

func getTokenCost(strategy *model.Strategy, tokenCount *big.Float) (*big.Float, error) {
	formTokenType := chain_service.NetworkInfoMap[strategy.NetWork].GasToken
	toTokenType := strategy.Token
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

func ValidateGas(userOp *model.UserOperation, gasComputeResponse *model.ComputeGasResponse, strategy *model.Strategy) error {
	paymasterDataExecutor := paymaster_pay_type.GetPaymasterDataExecutor(strategy.PayType)
	if paymasterDataExecutor == nil {
		return xerrors.Errorf(" %s paymasterDataExecutor not found", strategy.PayType)
	}
	return paymasterDataExecutor.ValidateGas(userOp, gasComputeResponse, strategy)
}
