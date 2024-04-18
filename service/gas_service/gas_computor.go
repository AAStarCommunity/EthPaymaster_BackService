package gas_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
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

	simulateResult, err := chain_service.SimulateHandleOp(strategy.GetNewWork())
	if err != nil {
		return nil, nil, err
	}

	feeResult := GetFeePerGas(strategy)
	preVerificationGas, err := chain_service.GetPreVerificationGas(strategy.GetNewWork(), userOp, strategy, feeResult)

	verificationGasLimit, err := EstimateVerificationGasLimit(strategy, simulateResult, preVerificationGas)

	callGasLimit, err := EstimateCallGasLimit(strategy, simulateResult, userOp)

	opEstimateGas := model.UserOpEstimateGas{}
	opEstimateGas.PreVerificationGas = preVerificationGas
	opEstimateGas.MaxFeePerGas = feeResult.MaxFeePerGas
	opEstimateGas.MaxPriorityFeePerGas = feeResult.MaxPriorityFeePerGas
	opEstimateGas.BaseFee = feeResult.BaseFee
	opEstimateGas.VerificationGasLimit = verificationGasLimit
	opEstimateGas.CallGasLimit = callGasLimit

	entryPointVersion := paymasterUserOp.GetEntrypointVersion()
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

	updateUserOp := GetNewUserOpAfterCompute(userOp, opEstimateGas)
	// TODO get PaymasterCallGasLimit
	return &model.ComputeGasResponse{
		GasInfo:       gasPrice,
		TokenCost:     tokenCost,
		OpEstimateGas: &opEstimateGas,
		Network:       strategy.GetNewWork(),
		Token:         strategy.GetUseToken(),
		UsdCost:       usdCost,
		MaxFee:        *maxFee,
	}, updateUserOp, nil
}

func GetNewUserOpAfterCompute(op *userop.BaseUserOp, gas model.UserOpEstimateGas) *userop.BaseUserOp {
	return nil
}

func GetFeePerGas(strategy *model.Strategy) (gasFeeResult *model.GasFeePerGasResult) {
	return nil
}

func EstimateCallGasLimit(strategy *model.Strategy, simulateOpResult *model.SimulateHandleOpResult, op *userop.BaseUserOp) (*big.Int, error) {
	ethereumExecutor := network.GetEthereumExecutor(strategy.GetNewWork())
	opValue := *op
	senderExist, _ := ethereumExecutor.CheckContractAddressAccess(opValue.GetSender())
	if senderExist {
		userOPCallGas, err := ethereumExecutor.EstimateUserOpCallGas(strategy.GetEntryPointAddress(), op)
		if err != nil {
			return nil, err
		}
		return userOPCallGas, nil
	} else {
		//1. TotalGas - createSenderGas = (verifyOpGas + verifyPaymasterGas) + callGasLimit
		//2. TotalGas -  (verifyOpGas + verifyPaymasterGas)  = executeUserOpGas；
		//3. executeUserOpGas（getFrom SimulateHandlop）- createSenderGas= callGasLimit
		initGas, err := ethereumExecutor.EstimateCreateSenderGas(strategy.GetEntryPointAddress(), op)
		if err != nil {
			return nil, err
		}
		executeUserOpGas := simulateOpResult.GasPaid
		return big.NewInt(0).Sub(executeUserOpGas, initGas), nil
	}
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

func EstimateVerificationGasLimit(strategy *model.Strategy, simulateOpResult *model.SimulateHandleOpResult, preVerificationGas *big.Int) (*big.Int, error) {
	preOpGas := simulateOpResult.PreOpGas
	result := new(big.Int).Sub(preOpGas, preVerificationGas)
	result = result.Mul(result, types.THREE_BIGINT)
	result = result.Div(result, types.TWO_BIGINT)
	if utils.LeftIsLessTanRight(result, types.DUMMY_VERIFICATIONGASLIMIT_BIGINT) {
		return types.DUMMY_VERIFICATIONGASLIMIT_BIGINT, nil
	}
	return result, nil
}
