package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/paymaster_data"
	"AAStarCommunity/EthPaymaster_BackService/common/price_compoent"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"AAStarCommunity/EthPaymaster_BackService/gas_executor"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"AAStarCommunity/EthPaymaster_BackService/service/validator_service"
	"AAStarCommunity/EthPaymaster_BackService/sponsor_manager"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"math/big"
	"strconv"
)

func TryPayUserOpExecute(apiKeyModel *model.ApiKeyModel, request *model.UserOpRequest) (*model.TryPayUserOpResponse, error) {
	userOp, strategy, paymasterDataInput, err := prepareExecute(request, apiKeyModel)
	if err != nil {
		return nil, err
	}

	gasResponse, paymasterUserOp, err := estimateGas(userOp, strategy, paymasterDataInput)
	if err != nil {
		return nil, err
	}

	paymasterDataInput.PaymasterVerificationGasLimit = gasResponse.OpEstimateGas.PaymasterVerificationGasLimit
	paymasterDataInput.PaymasterPostOpGasLimit = gasResponse.OpEstimateGas.PaymasterPostOpGasLimit

	//payReceipt, err := executePay(strategy, paymasterUserOp, gasResponse)
	//if err != nil {
	//	return nil, err
	//}
	//logrus.Debug("payReceipt:", payReceipt)
	result, err := postExecute(apiKeyModel, paymasterUserOp, strategy, gasResponse, paymasterDataInput)
	if err != nil {
		return nil, err
	}
	logrus.Debug("postExecute result:", result)
	return result, nil
}

func prepareExecute(request *model.UserOpRequest, apiKeyModel *model.ApiKeyModel) (*user_op.UserOpInput, *model.Strategy, *paymaster_data.PaymasterDataInput, error) {
	var strategy *model.Strategy
	strategy, generateErr := StrategyGenerate(request)
	if generateErr != nil {
		return nil, nil, nil, generateErr
	}

	if err := validator_service.ValidateStrategy(strategy, request, apiKeyModel); err != nil {
		return nil, nil, nil, err
	}

	userOp, err := user_op.NewUserOp(&request.UserOp)
	if err != nil {
		return nil, nil, nil, err
	}

	if err := validator_service.ValidateUserOp(userOp, strategy); err != nil {
		return nil, nil, nil, err
	}
	paymasterDataIput := paymaster_data.NewPaymasterDataInput(strategy)
	paymaster_data.NewPaymasterDataInput(strategy)
	return userOp, strategy, paymasterDataIput, nil
}

func estimateGas(userOp *user_op.UserOpInput, strategy *model.Strategy, paymasterDataInput *paymaster_data.PaymasterDataInput) (*model.ComputeGasResponse, *user_op.UserOpInput, error) {
	//base Strategy and UserOp computeGas
	gasResponse, paymasterUserOp, gasComputeError := gas_executor.ComputeGas(userOp, strategy, paymasterDataInput)
	if gasComputeError != nil {
		return nil, nil, gasComputeError
	}
	//The maxFeePerGas and maxPriorityFeePerGas are above a configurable minimum value that the client is willing to accept. At the minimum, they are sufficiently high to be included with the current block.basefee.
	if err := ValidateGas(userOp, gasResponse, strategy); err != nil {
		return nil, nil, err
	}
	return gasResponse, paymasterUserOp, nil
}

func ValidateGas(userOp *user_op.UserOpInput, gasComputeResponse *model.ComputeGasResponse, strategy *model.Strategy) error {
	validateFunc := gas_executor.GetGasValidateFunc(strategy.GetPayType())
	if validateFunc == nil {
		return xerrors.Errorf("ValidateGas Not Support PayType [%s]", strategy.GetPayType())
	}
	err := validateFunc(userOp, gasComputeResponse, strategy)
	if err != nil {
		return err
	}
	return nil
}

func executePay(input *ExecutePayInput) (*model.PayResponse, error) {
	if input.PayType == global_const.PayTypeERC20 {
		logrus.Debugf("Not Need ExecutePay In ERC20 PayType")
		return nil, nil
	}
	if config.IsSponsorWhitelist(input.UserOpSender) {
		logrus.Debugf("Not Need ExecutePay In SponsorWhitelist [%s]", input.UserOpSender)
		return nil, nil
	}
	if config.IsSponsorWhiteApiList(input.ApiKey) {
		logrus.Debugf("Not Need ExecutePay In SponsorWhiteKeylist [%s]", input.ApiKey)
		return nil, nil
	}

	// TODO
	//if config.IsTestNet(input.Network) {
	//	logrus.Debugf("Not Need ExecutePay In TestNet [%s]", input.Network)
	//	return nil, nil
	//}
	// Get Deposit Balance
	var payUserKey string
	if input.ProjectSponsor == true {
		payUserKey = input.ProjectUserId
	} else {
		payUserKey = input.UserOpSender
	}
	isTestNet := config.IsTestNet(input.Network)
	depositBalance, err := sponsor_manager.GetAvailableBalance(payUserKey, isTestNet)
	if err != nil {
		return nil, err
	}
	gasUsdCost, err := price_compoent.GetTokenCostInUsd(input.GasToken, input.MaxTxGasCostInEther)
	if err != nil {
		return nil, err
	}
	if depositBalance.Cmp(gasUsdCost) < 0 {
		return nil, xerrors.Errorf("Insufficient balance [%s] not Enough to Pay Cost [%s]", depositBalance.String(), gasUsdCost.String())
	}
	//Lock Deposit Balance
	_, err = sponsor_manager.LockUserBalance(payUserKey, input.UserOpHash, isTestNet,
		gasUsdCost)
	if err != nil {
		return nil, err
	}

	return &model.PayResponse{
		PayType: input.PayType,
	}, nil
}

type ExecutePayInput struct {
	ApiKey              string
	ProjectUserId       string
	PayType             global_const.PayType
	ProjectSponsor      bool
	UserOpSender        string
	MaxTxGasCostInEther *big.Float
	UserOpHash          []byte
	Network             global_const.Network
	GasToken            global_const.TokenType
}

func postExecute(apiKeyModel *model.ApiKeyModel, userOp *user_op.UserOpInput, strategy *model.Strategy, gasResponse *model.ComputeGasResponse, paymasterDataInput *paymaster_data.PaymasterDataInput) (*model.TryPayUserOpResponse, error) {

	executor := network.GetEthereumExecutor(strategy.GetNewWork())
	paymasterData, userOpHash, err := executor.GetPaymasterData(userOp, strategy, paymasterDataInput)
	if err != nil {
		return nil, xerrors.Errorf("postExecute GetPaymasterData Error: [%w]", err)
	}
	logrus.Debug("postExecute paymasterData:", paymasterData)

	_, err = executePay(&ExecutePayInput{
		ApiKey:              apiKeyModel.ApiKey,
		ProjectUserId:       strconv.FormatInt(apiKeyModel.UserId, 10),
		PayType:             strategy.GetPayType(),
		ProjectSponsor:      strategy.ProjectSponsor,
		UserOpSender:        userOp.Sender.String(),
		MaxTxGasCostInEther: gasResponse.TotalGasDetail.MaxTxGasCostInEther,
		UserOpHash:          userOpHash,
		Network:             strategy.GetNewWork(),
		GasToken:            strategy.GetGasToken(),
	})
	if err != nil {
		return nil, xerrors.Errorf("postExecute executePay Error: [%w]", err)
	}
	var result = &model.TryPayUserOpResponse{
		StrategyId:        strategy.Id,
		EntryPointAddress: strategy.GetEntryPointAddress().String(),
		NetWork:           strategy.GetNewWork(),
		EntrypointVersion: strategy.GetStrategyEntrypointVersion(),
		PayMasterAddress:  strategy.GetPaymasterAddress().String(),
		Erc20TokenCost:    gasResponse.Erc20TokenCost,

		UserOpResponse: &model.UserOpResponse{
			PayMasterAndData:     utils.EncodeToHexStringWithPrefix(paymasterData),
			PreVerificationGas:   utils.ConvertBigIntToHexWithPrefix(gasResponse.OpEstimateGas.PreVerificationGas),
			MaxFeePerGas:         utils.ConvertBigIntToHexWithPrefix(gasResponse.OpEstimateGas.MaxFeePerGas),
			MaxPriorityFeePerGas: utils.ConvertBigIntToHexWithPrefix(gasResponse.OpEstimateGas.MaxPriorityFeePerGas),
			VerificationGasLimit: utils.ConvertBigIntToHexWithPrefix(gasResponse.OpEstimateGas.VerificationGasLimit),
			CallGasLimit:         utils.ConvertBigIntToHexWithPrefix(gasResponse.OpEstimateGas.CallGasLimit),
		},
	}

	if strategy.GetStrategyEntrypointVersion() == global_const.EntrypointV07 {
		result.UserOpResponse.AccountGasLimit = utils.EncodeToHexStringWithPrefix(gasResponse.OpEstimateGas.AccountGasLimit[:])
		result.UserOpResponse.GasFees = utils.EncodeToHexStringWithPrefix(gasResponse.OpEstimateGas.GasFees[:])
		result.UserOpResponse.PaymasterVerificationGasLimit = utils.ConvertBigIntToHexWithPrefix(gasResponse.OpEstimateGas.PaymasterVerificationGasLimit)
		result.UserOpResponse.PaymasterPostOpGasLimit = utils.ConvertBigIntToHexWithPrefix(gasResponse.OpEstimateGas.PaymasterPostOpGasLimit)
	}

	return result, nil
}

func StrategyGenerate(request *model.UserOpRequest) (*model.Strategy, error) {
	var strategyResult *model.Strategy
	if forceStrategyId := request.StrategyCode; forceStrategyId != "" {
		//force strategy
		strategy, err := dashboard_service.GetStrategyByCode(forceStrategyId, request.EntryPointVersion, request.Network)
		if err != nil {
			return nil, err
		}
		if strategy == nil {
			return nil, xerrors.Errorf("Empty Strategies")
		}
		strategyResult = strategy

	} else {
		suitableStrategy, err := dashboard_service.GetSuitableStrategyWithOutCode(request.EntryPointVersion, request.Network, request.UserPayErc20Token)
		if err != nil {
			return nil, err
		}
		if suitableStrategy == nil {
			return nil, xerrors.Errorf("Empty Strategies")
		}

		strategyResult = suitableStrategy
	}
	return strategyResult, nil
}
