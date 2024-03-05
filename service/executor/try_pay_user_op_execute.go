package executor

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"AAStarCommunity/EthPaymaster_BackService/service/gas_service"
	"AAStarCommunity/EthPaymaster_BackService/service/validator_service"
	"golang.org/x/xerrors"
)

func TryPayUserOpExecute(request model.TryPayUserOpRequest) (model.Result, error) {
	//validator
	if err := paramValidate(request); err != nil {
		return model.Result{}, err
	}
	userOp := request.UserOperation
	//getStrategy
	strategy, err := strategyGenerate(request)
	if err != nil {
		return model.Result{}, err
	}
	if err := validator_service.ValidateStrategy(strategy, userOp); err != nil {
		return model.Result{}, err
	}

	//base Strategy and UserOp computeGas
	gasResponse, gasComputeError := gas_service.ComputeGas(userOp, strategy)
	if gasComputeError != nil {
		return model.Result{}, gasComputeError
	}

	//validate gas
	if err := gas_service.ValidateGas(userOp, gasResponse); err != nil {
		return model.Result{}, err
	}
	//pay
	payReceipt, payError := executePay(strategy, request.UserOperation, gasResponse)
	if payError != nil {
		return model.Result{}, payError
	}
	paymasterSignature := getPayMasterSignature(strategy, request.UserOperation)
	var result = model.TryPayUserOpResponse{
		StrategyId:         strategy.Id,
		EntryPointAddress:  strategy.EntryPointAddress,
		PayMasterAddress:   strategy.EntryPointAddress,
		PayReceipt:         payReceipt,
		PayMasterSignature: paymasterSignature,
		GasInfo:            gasResponse,
	}

	return model.Result{
		Code:    0,
		Data:    result,
		Message: "message",
		Cost:    "cost",
	}, nil
}
func paramValidate(request model.TryPayUserOpRequest) error {
	return nil
}

func executePay(strategy model.Strategy, userOp model.UserOperationItem, gasResponse model.ComputeGasResponse) (interface{}, error) {
	//1.Recharge
	//2.record account
	//3.return Receipt
	return nil, nil
}
func getPayMasterSignature(strategy model.Strategy, userOp model.UserOperationItem) string {
	return ""
}

func strategyGenerate(request model.TryPayUserOpRequest) (model.Strategy, error) {
	if forceStrategyId := request.ForceStrategyId; forceStrategyId != "" {
		//force strategy
		strategy := dashboard_service.GetStrategyById(forceStrategyId)
		if strategy == (model.Strategy{}) {
			return model.Strategy{}, xerrors.Errorf("Not Support Strategy ID: [%w]", forceStrategyId)
		}
		return strategy, nil
	}

	suitableStrategy, err := dashboard_service.GetSuitableStrategy(request.ForceEntryPointAddress, request.ForceNetWork, request.ForceTokens) //TODO
	if err != nil {
		return model.Strategy{}, err
	}
	if suitableStrategy == (model.Strategy{}) {
		return model.Strategy{}, xerrors.Errorf("Empty Strategies")
	}
	return suitableStrategy, nil
}
