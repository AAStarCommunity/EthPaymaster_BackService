package gas_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/data_utils"
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/paymaster_data"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"AAStarCommunity/EthPaymaster_BackService/gas_validate"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"math/big"
)

var (
	GweiFactor   = new(big.Float).SetInt(big.NewInt(1e9))
	EthWeiFactor = new(big.Float).SetInt(big.NewInt(1e18))
)

// https://blog.particle.network/bundler-predicting-gas/
func ComputeGas(userOp *user_op.UserOpInput, strategy *model.Strategy, paymasterDataInput *paymaster_data.PaymasterData) (*model.ComputeGasResponse, *user_op.UserOpInput, error) {

	opEstimateGas, err := getUserOpEstimateGas(userOp, strategy, paymasterDataInput)
	if err != nil {
		return nil, nil, err
	}

	totalGasDetail := GetTotalCostByEstimateGas(opEstimateGas)
	updateUserOp := getNewUserOpAfterCompute(userOp, opEstimateGas, strategy.GetStrategyEntrypointVersion())
	return &model.ComputeGasResponse{
		OpEstimateGas:  opEstimateGas,
		TotalGasDetail: totalGasDetail,
	}, updateUserOp, nil
}
func GetTotalCostByEstimateGas(userOpGas *model.UserOpEstimateGas) *model.TotalGasDetail {
	gasPrice := GetUserOpGasPrice(userOpGas)
	totalGasLimit := new(big.Int)
	totalGasLimit = totalGasLimit.Add(totalGasLimit, userOpGas.VerificationGasLimit)
	totalGasLimit = totalGasLimit.Add(totalGasLimit, userOpGas.CallGasLimit)
	totalGasLimit = totalGasLimit.Add(totalGasLimit, userOpGas.PreVerificationGas)
	totalGasGost := new(big.Int).Mul(gasPrice, totalGasLimit)

	gasPriceInGwei := new(big.Float).SetInt(gasPrice)
	gasPriceInGwei.Quo(gasPriceInGwei, GweiFactor)

	totalGasGostInGwei := new(big.Float).SetInt(totalGasGost)
	totalGasGostInGwei.Quo(totalGasGostInGwei, GweiFactor)
	logrus.Debug("totalGasGostInGwei: ", totalGasGostInGwei)

	totalGasGostInEther := new(big.Float).SetInt(totalGasGost)
	totalGasGostInEther.Quo(totalGasGostInEther, EthWeiFactor)
	logrus.Debug("totalGasGostInEther: ", totalGasGostInEther)

	return &model.TotalGasDetail{
		MaxTxGasLimit:       totalGasLimit,
		MaxTxGasCostGwei:    totalGasGostInGwei,
		MaxTxGasCostInEther: totalGasGostInEther,
		GasPriceGwei:        gasPriceInGwei,
	}
}

// GetUserOpGasPrice if network not Support EIP1559 will set MaxFeePerGas And  MaxPriorityFeePerGas to the same value
func GetUserOpGasPrice(userOpGas *model.UserOpEstimateGas) *big.Int {
	maxFeePerGas := userOpGas.MaxFeePerGas
	maxPriorityFeePerGas := userOpGas.MaxPriorityFeePerGas
	if maxFeePerGas == maxPriorityFeePerGas {
		return maxFeePerGas
	}
	combineFee := new(big.Int).Add(userOpGas.BaseFee, maxPriorityFeePerGas)
	return utils.GetMinValue(maxFeePerGas, combineFee)
}

func getUserOpEstimateGas(userOp *user_op.UserOpInput, strategy *model.Strategy, paymasterDataInput *paymaster_data.PaymasterData) (*model.UserOpEstimateGas, error) {
	gasPriceResult, gasPriceErr := chain_service.GetGasPrice(strategy.GetNewWork())
	if userOp.MaxFeePerGas != nil {
		gasPriceResult.MaxFeePerGas = userOp.MaxFeePerGas
	}
	if userOp.MaxPriorityFeePerGas != nil {
		gasPriceResult.MaxPriorityFeePerGas = userOp.MaxPriorityFeePerGas
	}

	//TODO calculate the maximum possible fee the account needs to pay (based on validation and call gas limits, and current gas values)
	if gasPriceErr != nil {
		return nil, xerrors.Errorf("get gas price error: %v", gasPriceErr)
	}
	userOpInputForSimulate, err := data_utils.GetUserOpWithPaymasterAndDataForSimulate(*userOp, strategy, paymasterDataInput, gasPriceResult)
	simulateGasPrice := utils.GetGasEntryPointGasGrace(gasPriceResult.MaxFeePerGas, gasPriceResult.MaxPriorityFeePerGas, gasPriceResult.BaseFee)
	if err != nil {
		return nil, xerrors.Errorf("GetUserOpWithPaymasterAndDataForSimulate error: %v", err)
	}

	simulateResult, err := chain_service.SimulateHandleOp(userOpInputForSimulate, strategy)
	if err != nil {
		return nil, xerrors.Errorf("SimulateHandleOp error: %v", err)
	}

	preVerificationGas, err := chain_service.GetPreVerificationGas(userOp, strategy, gasPriceResult, simulateResult)

	verificationGasLimit, err := estimateVerificationGasLimit(simulateResult, preVerificationGas)

	callGasLimit, err := EstimateCallGasLimit(strategy, simulateResult, userOp, simulateGasPrice)

	opEstimateGas := model.UserOpEstimateGas{}
	opEstimateGas.PreVerificationGas = preVerificationGas
	opEstimateGas.MaxFeePerGas = gasPriceResult.MaxFeePerGas
	opEstimateGas.MaxPriorityFeePerGas = gasPriceResult.MaxPriorityFeePerGas
	opEstimateGas.BaseFee = gasPriceResult.BaseFee
	opEstimateGas.VerificationGasLimit = verificationGasLimit
	opEstimateGas.CallGasLimit = callGasLimit

	entryPointVersion := strategy.GetStrategyEntrypointVersion()
	if entryPointVersion == global_const.EntryPointV07 {
		opEstimateGas.AccountGasLimit = utils.PackIntTo32Bytes(verificationGasLimit, callGasLimit)
		opEstimateGas.GasFees = utils.PackIntTo32Bytes(gasPriceResult.MaxPriorityFeePerGas, gasPriceResult.MaxFeePerGas)
		opEstimateGas.PaymasterPostOpGasLimit = global_const.DummyPaymasterPostopGaslimitBigint
		opEstimateGas.PaymasterVerificationGasLimit = global_const.DummyPaymasterVerificationgaslimitBigint
	}
	return &opEstimateGas, nil
}

func getNewUserOpAfterCompute(op *user_op.UserOpInput, gas *model.UserOpEstimateGas, version global_const.EntrypointVersion) *user_op.UserOpInput {
	var accountGasLimits [32]byte
	var gasFee [32]byte
	if version == global_const.EntryPointV07 {
		accountGasLimits = utils.PackIntTo32Bytes(gas.PreVerificationGas, gas.CallGasLimit)
		gasFee = utils.PackIntTo32Bytes(gas.MaxPriorityFeePerGas, gas.MaxFeePerGas)
	}
	result := &user_op.UserOpInput{
		Sender:               op.Sender,
		Nonce:                op.Nonce,
		InitCode:             op.InitCode,
		CallData:             op.CallData,
		MaxFeePerGas:         op.MaxFeePerGas,
		Signature:            op.Signature,
		MaxPriorityFeePerGas: op.MaxPriorityFeePerGas,
		CallGasLimit:         op.CallGasLimit,
		VerificationGasLimit: op.VerificationGasLimit,
		AccountGasLimits:     accountGasLimits,
		GasFees:              gasFee,
		PreVerificationGas:   op.PreVerificationGas,
	}
	return result
}

func EstimateCallGasLimit(strategy *model.Strategy, simulateOpResult *model.SimulateHandleOpResult, op *user_op.UserOpInput, simulateGasPrice *big.Int) (*big.Int, error) {
	ethereumExecutor := network.GetEthereumExecutor(strategy.GetNewWork())
	opValue := *op
	senderExist, _ := ethereumExecutor.CheckContractAddressAccess(opValue.Sender)
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
		executeUserOpGas := new(big.Int).Div(simulateOpResult.GasPaid, simulateGasPrice)
		return big.NewInt(0).Sub(executeUserOpGas, initGas), nil
	}
}

func getTokenCost(strategy *model.Strategy, tokenCount *big.Float) (*big.Float, error) {
	if strategy.GetPayType() == global_const.PayTypeERC20 {

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

func ValidateGas(userOp *user_op.UserOpInput, gasComputeResponse *model.ComputeGasResponse, strategy *model.Strategy) error {
	validateFunc := gas_validate.GasValidateFuncMap[strategy.GetPayType()]
	err := validateFunc(userOp, gasComputeResponse, strategy)
	if err != nil {
		return err
	}
	return nil
}

func estimateVerificationGasLimit(simulateOpResult *model.SimulateHandleOpResult, preVerificationGas *big.Int) (*big.Int, error) {
	preOpGas := simulateOpResult.PreOpGas
	// verificationGasLimit = (preOpGas - preVerificationGas) * 1.5
	result := new(big.Int).Sub(preOpGas, preVerificationGas)
	result = result.Mul(result, global_const.ThreeBigint)
	result = result.Div(result, global_const.TwoBigint)
	if utils.LeftIsLessTanRight(result, global_const.DummyVerificationgaslimitBigint) {
		return global_const.DummyVerificationgaslimitBigint, nil
	}
	return result, nil
}
