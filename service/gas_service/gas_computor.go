package gas_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
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
	maxFeePriceInEther := new(big.Float).SetInt(maxFee)
	maxFeePriceInEther.Quo(maxFeePriceInEther, chain_service.EthWeiFactor)
	tokenCost, _ := getTokenCost(strategy, maxFeePriceInEther)
	if strategy.PayType == types.PayTypeERC20 {
		//TODO get ERC20 balance
		if err := validateErc20Paymaster(tokenCost, strategy); err != nil {
			return nil, err
		}
	}

	// TODO get PaymasterCallGasLimit
	return &model.ComputeGasResponse{
		GasInfo:    gasPrice,
		TokenCost:  tokenCost.Text('f', 18),
		Network:    strategy.NetWork,
		Token:      strategy.Token,
		UsdCost:    "0.4",
		BlobEnable: strategy.Enable4844,
		MaxFee:     *maxFee,
	}, nil
}
func validateErc20Paymaster(tokenCost *big.Float, strategy *model.Strategy) error {
	//useToken := strategy.Token
	//// get User address balance
	//TODO
	return nil
}
func getTokenCost(strategy *model.Strategy, tokenCount *big.Float) (*big.Float, error) {
	formTokenType := chain_service.NetworkInfoMap[strategy.NetWork].GasToken
	toTokenType := strategy.Token
	toTokenPrice, err := utils.GetToken(formTokenType, toTokenType)
	if err != nil {
		return nil, err
	}
	tokenCost := new(big.Float).Mul(tokenCount, big.NewFloat(toTokenPrice))
	return tokenCost, nil
}
func GetPayMasterGasLimit() *big.Int {
	return nil
}

func ValidateGas(userOp *model.UserOperation, gasComputeResponse *model.ComputeGasResponse, strategy *model.Strategy) error {
	//1.if ERC20 check address balacnce
	//Validate the account’s deposit in the entryPoint is high enough to cover the max possible cost (cover the already-done verification and max execution gas)
	//2 if Paymaster check paymaster balance
	//The maxFeePerGas and maxPriorityFeePerGas are above a configurable minimum value that the client is willing to accept. At the minimum, they are sufficiently high to be included with the current block.basefee.
	if strategy.PayType == types.PayTypeERC20 {
		//TODO check address balance
	} else if strategy.PayType == types.PayTypeVerifying {
		//TODO check paymaster balance
	}
	return nil
}
