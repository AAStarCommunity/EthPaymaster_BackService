package gas_executor

import (
	"AAStarCommunity/EthPaymaster_BackService/common/data_utils"
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/paymaster_data"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"math/big"
)

var (
	GweiFactor   = new(big.Float).SetInt(big.NewInt(1e9))
	EthWeiFactor = new(big.Float).SetInt(big.NewInt(1e18))
)

// ComputeGas https://blog.particle.network/bundler-predicting-gas/
func ComputeGas(userOp *user_op.UserOpInput, strategy *model.Strategy, paymasterDataInput *paymaster_data.PaymasterDataInput) (*model.ComputeGasResponse, *user_op.UserOpInput, error) {

	opEstimateGas, err := getUserOpEstimateGas(userOp, strategy, paymasterDataInput)
	if err != nil {
		return nil, nil, err
	}

	totalGasDetail := GetTotalCostByEstimateGas(opEstimateGas)
	updateUserOp := getNewUserOpAfterCompute(userOp, opEstimateGas, strategy.GetStrategyEntrypointVersion())
	var erc20TokenCost *big.Float
	if !userOp.ComputeGasOnly {
		erc20TokenCost, err = getErc20TokenCost(strategy, totalGasDetail.MaxTxGasCostInEther)
		if err != nil {
			return nil, nil, xerrors.Errorf("getErc20TokenCost error: %v", err)
		}
	}
	return &model.ComputeGasResponse{
		OpEstimateGas:  opEstimateGas,
		TotalGasDetail: totalGasDetail,
		Erc20TokenCost: erc20TokenCost,
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

func getUserOpEstimateGas(userOp *user_op.UserOpInput, strategy *model.Strategy, paymasterDataInput *paymaster_data.PaymasterDataInput) (*model.UserOpEstimateGas, error) {
	gasPriceResult, gasPriceErr := GetGasPrice(strategy.GetNewWork())
	if gasPriceErr != nil {
		return nil, xerrors.Errorf("get gas price error: %v", gasPriceErr)
	}
	if userOp.MaxFeePerGas != nil && userOp.MaxPriorityFeePerGas != nil {
		if config.IsDisable1559Chain(strategy.GetNewWork()) && userOp.MaxFeePerGas.Cmp(userOp.MaxPriorityFeePerGas) != 0 {
			return nil, xerrors.Errorf("[%v] is not support 1559 MaxFeePerGas and MaxPriorityFeePerGas can not be same  at the same time", strategy.GetNewWork())
		}
	}
	if userOp.MaxFeePerGas != nil {
		gasPriceResult.MaxFeePerGas = userOp.MaxFeePerGas
	}
	if userOp.MaxPriorityFeePerGas != nil {
		gasPriceResult.MaxPriorityFeePerGas = userOp.MaxPriorityFeePerGas
	}

	//TODO calculate the maximum possible fee the account needs to pay (based on validation and call gas limits, and current gas values)

	userOpInputForSimulate, err := data_utils.GetUserOpWithPaymasterAndDataForSimulate(*userOp, strategy, paymasterDataInput, gasPriceResult)
	if err != nil {
		return nil, xerrors.Errorf("GetUserOpWithPaymasterAndDataForSimulate error: %v", err)
	}
	logrus.Debugf("userOpInputForSimulate: %v", userOpInputForSimulate)
	logrus.Debugf("getUserOpEstimateGas gasPriceResult: %v", gasPriceResult)
	simulateGasPrice := utils.GetGasEntryPointGasGrace(gasPriceResult.MaxFeePerGas, gasPriceResult.MaxPriorityFeePerGas, gasPriceResult.BaseFee)

	simulateResult, err := chain_service.SimulateHandleOp(userOpInputForSimulate, strategy)
	if err != nil {
		return nil, xerrors.Errorf("SimulateHandleOp error: %v", err)
	}
	preVerificationGas, err := GetPreVerificationGas(userOp, strategy, gasPriceResult, simulateResult)
	if err != nil {
		return nil, xerrors.Errorf("GetPreVerificationGas error: %v", err)
	}
	verificationGasLimit, err := estimateVerificationGasLimit(simulateResult, preVerificationGas)
	if err != nil {
		return nil, xerrors.Errorf("estimateVerificationGasLimit error: %v", err)
	}

	callGasLimit, err := EstimateCallGasLimit(strategy, simulateResult, userOp, simulateGasPrice)
	if err != nil {
		return nil, xerrors.Errorf("EstimateCallGasLimit error: %v", err)
	}

	opEstimateGas := model.UserOpEstimateGas{}
	opEstimateGas.PreVerificationGas = preVerificationGas
	opEstimateGas.MaxFeePerGas = gasPriceResult.MaxFeePerGas
	opEstimateGas.MaxPriorityFeePerGas = gasPriceResult.MaxPriorityFeePerGas
	opEstimateGas.BaseFee = gasPriceResult.BaseFee
	opEstimateGas.VerificationGasLimit = verificationGasLimit
	opEstimateGas.CallGasLimit = callGasLimit

	entryPointVersion := strategy.GetStrategyEntrypointVersion()
	if entryPointVersion == global_const.EntrypointV07 {
		accountGasLimit := utils.PackIntTo32Bytes(verificationGasLimit, callGasLimit)
		opEstimateGas.AccountGasLimit = &accountGasLimit
		gasFees := utils.PackIntTo32Bytes(gasPriceResult.MaxPriorityFeePerGas, gasPriceResult.MaxFeePerGas)
		opEstimateGas.GasFees = &gasFees
		opEstimateGas.PaymasterPostOpGasLimit = global_const.DummyPaymasterPostoperativelyBigint
		opEstimateGas.PaymasterVerificationGasLimit = global_const.DummyPaymasterOversimplificationBigint
	}
	return &opEstimateGas, nil
}

func getNewUserOpAfterCompute(op *user_op.UserOpInput, gas *model.UserOpEstimateGas, version global_const.EntrypointVersion) *user_op.UserOpInput {
	var accountGasLimits [32]byte
	var gasFee [32]byte
	if version == global_const.EntrypointV07 {
		accountGasLimits = utils.PackIntTo32Bytes(gas.PreVerificationGas, gas.CallGasLimit)
		gasFee = utils.PackIntTo32Bytes(gas.MaxPriorityFeePerGas, gas.MaxFeePerGas)
	}
	result := &user_op.UserOpInput{
		Sender:               op.Sender,
		Nonce:                op.Nonce,
		InitCode:             op.InitCode,
		CallData:             op.CallData,
		MaxFeePerGas:         gas.MaxFeePerGas,
		Signature:            op.Signature,
		MaxPriorityFeePerGas: gas.MaxPriorityFeePerGas,
		CallGasLimit:         gas.CallGasLimit,
		VerificationGasLimit: gas.VerificationGasLimit,
		AccountGasLimits:     accountGasLimits,
		GasFees:              gasFee,
		PreVerificationGas:   gas.PreVerificationGas,
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

func getErc20TokenCost(strategy *model.Strategy, tokenCount *big.Float) (*big.Float, error) {
	if strategy.GetPayType() == global_const.PayTypeERC20 {
		if strategy.Erc20TokenType == "" {
			return nil, xerrors.Errorf("strategy.Erc20TokenType is nil")
		}
		formTokenType := config.GetGasToken(strategy.GetNewWork())
		toTokenType := strategy.Erc20TokenType
		toTokenPrice, err := utils.GetToken(formTokenType, toTokenType)
		if err != nil {
			return nil, err
		}
		if toTokenPrice == 0 {
			return nil, xerrors.Errorf("toTokenPrice can not be 0")
		}
		logrus.Debugf("toTokenPrice: %v", toTokenPrice)
		tokenCost := new(big.Float).Mul(tokenCount, big.NewFloat(toTokenPrice))
		return tokenCost, nil
	}
	return nil, nil

}

func estimateVerificationGasLimit(simulateOpResult *model.SimulateHandleOpResult, preVerificationGas *big.Int) (*big.Int, error) {
	preOpGas := simulateOpResult.PreOpGas
	logrus.Debugf("preOpGas: %v", preOpGas)
	logrus.Debugf("preVerificationGas: %v", preVerificationGas)
	// verificationGasLimit = (preOpGas - preVerificationGas) * 1.5
	result := new(big.Int).Sub(preOpGas, preVerificationGas)
	result = result.Mul(result, global_const.ThreeBigint)
	result = result.Div(result, global_const.TwoBigint)
	if utils.LeftIsLessTanRight(result, global_const.DummyVerificationGasLimit) {
		return global_const.DummyVerificationGasLimit, nil
	}
	return result, nil
}

// GetGasPrice return gas price in wei, gwei, ether
func GetGasPrice(chain global_const.Network) (*model.GasPrice, error) {
	if config.IsEthereumAdaptableNetWork(chain) {
		ethereumExecutor := network.GetEthereumExecutor(chain)
		return ethereumExecutor.GetGasPrice()
	} else if chain == global_const.StarketMainnet || chain == global_const.StarketSepolia {
		starknetExecutor := network.GetStarknetExecutor()
		return starknetExecutor.GetGasPrice()
	} else {
		return nil, xerrors.Errorf("chain %s not support", chain)
	}
	//MaxFeePerGas
	//MaxPriorityPrice
	//preOpGas (get verificationGasLimit from preOpGas)
	//

}

// GetPreVerificationGas https://github.com/eth-infinitism/bundler/blob/main/packages/sdk/src/calcPreVerificationGas.ts
func GetPreVerificationGas(userOp *user_op.UserOpInput, strategy *model.Strategy, gasFeeResult *model.GasPrice, simulateOpResult *model.SimulateHandleOpResult) (*big.Int, error) {
	chain := strategy.GetNewWork()
	stack := config.GetNetWorkStack(chain)
	preGasFunc, err := GetPreVerificationGasFunc(stack)
	if err != nil {
		return nil, err
	}
	preGas, err := preGasFunc(&model.PreVerificationGasEstimateInput{
		Strategy:         strategy,
		Op:               userOp,
		GasFeeResult:     gasFeeResult,
		SimulateOpResult: simulateOpResult,
	})
	if err != nil {
		return nil, err
	}
	// add 10% buffer
	preGas = preGas.Mul(preGas, global_const.HundredPlusOneBigint)
	preGas = preGas.Div(preGas, global_const.HundredBigint)
	logrus.Debugf("GetPreVerificationGas preVerificationGas: %v", preGas)
	return preGas, nil
}
