package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/paymaster_data"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"AAStarCommunity/EthPaymaster_BackService/service/gas_service"
	"AAStarCommunity/EthPaymaster_BackService/service/pay_service"
	"AAStarCommunity/EthPaymaster_BackService/service/validator_service"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

func TryPayUserOpExecute(request *model.UserOpRequest) (*model.TryPayUserOpResponse, error) {
	userOp, strategy, paymasterDtataIput, err := prepareExecute(request)
	if err != nil {
		return nil, err
	}

	gasResponse, paymasterUserOp, err := estimateGas(userOp, strategy, paymasterDtataIput)
	if err != nil {
		return nil, err
	}

	payReceipt, err := executePay(strategy, paymasterUserOp, gasResponse)
	if err != nil {
		return nil, err
	}
	logrus.Debug("payReceipt:", payReceipt)
	result, err := postExecute(paymasterUserOp, strategy, gasResponse, paymasterDtataIput)
	if err != nil {
		return nil, err
	}
	logrus.Debug("postExecute result:", result)
	result.PayReceipt = payReceipt
	return result, nil
}

//sub Function ---------

func prepareExecute(request *model.UserOpRequest) (*user_op.UserOpInput, *model.Strategy, *paymaster_data.PaymasterData, error) {

	var strategy *model.Strategy

	strategy, generateErr := StrategyGenerate(request)
	if generateErr != nil {
		return nil, nil, nil, generateErr
	}

	userOp, err := user_op.NewUserOp(&request.UserOp)
	if err != nil {
		return nil, nil, nil, err

	}
	if err := validator_service.ValidateStrategy(strategy); err != nil {
		return nil, nil, nil, err
	}
	if err := validator_service.ValidateUserOp(userOp, strategy); err != nil {
		return nil, nil, nil, err
	}
	paymasterDataIput := paymaster_data.NewPaymasterDataInput(strategy)
	paymaster_data.NewPaymasterDataInput(strategy)
	return userOp, strategy, paymasterDataIput, nil
}

func estimateGas(userOp *user_op.UserOpInput, strategy *model.Strategy, paymasterDataInput *paymaster_data.PaymasterData) (*model.ComputeGasResponse, *user_op.UserOpInput, error) {
	//base Strategy and UserOp computeGas
	gasResponse, paymasterUserOp, gasComputeError := gas_service.ComputeGas(userOp, strategy, paymasterDataInput)
	if gasComputeError != nil {
		return nil, nil, gasComputeError
	}
	//The maxFeePerGas and maxPriorityFeePerGas are above a configurable minimum value that the client is willing to accept. At the minimum, they are sufficiently high to be included with the current block.basefee.
	//validate gas
	if err := gas_service.ValidateGas(userOp, gasResponse, strategy); err != nil {
		return nil, nil, err
	}
	return gasResponse, paymasterUserOp, nil
}

func executePay(strategy *model.Strategy, userOp *user_op.UserOpInput, gasResponse *model.ComputeGasResponse) (*model.PayReceipt, error) {
	//1.Recharge
	ethereumPayservice := pay_service.EthereumPayService{}
	if err := ethereumPayservice.Pay(); err != nil {
		return nil, err
	}
	//2.record account
	ethereumPayservice.RecordAccount()
	//3.return Receipt
	ethereumPayservice.GetReceipt()
	return &model.PayReceipt{
		TransactionHash: "0x110406d44ec1681fcdab1df2310181dee26ff43c37167b2c9c496b35cce69437",
		Sponsor:         "aastar",
	}, nil
}

func postExecute(userOp *user_op.UserOpInput, strategy *model.Strategy, gasResponse *model.ComputeGasResponse, paymasterDataInput *paymaster_data.PaymasterData) (*model.TryPayUserOpResponse, error) {
	executor := network.GetEthereumExecutor(strategy.GetNewWork())
	paymasterData, err := executor.GetPaymasterData(userOp, strategy, paymasterDataInput)
	if err != nil {
		return nil, xerrors.Errorf("postExecute GetPaymasterData Error: [%w]", err)
	}
	logrus.Debug("postExecute paymasterData:", paymasterData)
	var result = &model.TryPayUserOpResponse{
		StrategyId:        strategy.Id,
		EntryPointAddress: strategy.GetEntryPointAddress().String(),
		PayMasterAddress:  strategy.GetPaymasterAddress().String(),
		PayMasterAndData:  utils.EncodeToStringWithPrefix(paymasterData),
		GasInfo:           gasResponse,
	}
	return result, nil
}

func StrategyGenerate(request *model.UserOpRequest) (*model.Strategy, error) {
	if forceStrategyId := request.ForceStrategyId; forceStrategyId != "" {
		//force strategy
		if strategy := dashboard_service.GetStrategyById(forceStrategyId); strategy == nil {
			return nil, xerrors.Errorf("Not Support Strategy ID: [%w]", forceStrategyId)
		} else {
			return strategy, nil
		}
	}

	suitableStrategy, err := dashboard_service.GetSuitableStrategy(request.ForceEntryPointAddress, request.ForceNetwork, global_const.PayTypeSuperVerifying) //TODO
	if err != nil {
		return nil, err
	}
	if suitableStrategy == nil {
		return nil, xerrors.Errorf("Empty Strategies")
	}
	return suitableStrategy, nil
}
