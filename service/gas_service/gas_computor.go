package gas_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
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
	//x := gasPrice.MaxBasePriceWei.Int64() + gasPrice.MaxPriorityPriceWei.Int64()
	//maxFeePerGas := (x, userOp.MaxFeePerGas.Uint64())
	payMasterPostGasLimit := GetPayMasterGasLimit()

	maxGasLimit := big.NewInt(0).Add(userOp.CallGasLimit, userOp.VerificationGasLimit)
	maxGasLimit = maxGasLimit.Add(maxGasLimit, payMasterPostGasLimit)

	maxFee := new(big.Int).Mul(maxGasLimit, gasPrice.MaxBasePriceWei)
	// TODO get PaymasterCallGasLimit
	tokenCost := GetTokenCost(*maxFee, userOp, *strategy)
	return &model.ComputeGasResponse{
		GasInfo:    gasPrice,
		TokenCost:  tokenCost,
		Network:    strategy.NetWork,
		Token:      strategy.Token,
		UsdCost:    "0.4",
		BlobEnable: strategy.Enable4844,
		MaxFee:     *maxFee,
	}, nil
}
func GetPayMasterGasLimit() *big.Int {
	return nil
}
func GetTokenCost(maxFee big.Int, userOp *model.UserOperation, strategy model.Strategy) string {
	return "0.0001"
}

func ValidateGas(userOp *model.UserOperation, gasComputeResponse *model.ComputeGasResponse) error {
	//1.if ERC20 check address balacnce
	//Validate the account’s deposit in the entryPoint is high enough to cover the max possible cost (cover the already-done verification and max execution gas)
	//2 if Paymaster check paymaster balance
	//The maxFeePerGas and maxPriorityFeePerGas are above a configurable minimum value that the client is willing to accept. At the minimum, they are sufficiently high to be included with the current block.basefee.
	return nil
}
