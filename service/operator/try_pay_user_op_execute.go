package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"AAStarCommunity/EthPaymaster_BackService/service/gas_service"
	"AAStarCommunity/EthPaymaster_BackService/service/validator_service"
	"golang.org/x/xerrors"
)

func TryPayUserOpExecute(request *model.TryPayUserOpRequest) (*model.Result, error) {
	// validator
	if err := businessParamValidate(request); err != nil {
		return nil, err
	}
	userOp := request.UserOperation

	// getStrategy
	var strategy *model.Strategy
	if stg, err := strategyGenerate(request); err != nil {
		return nil, err
	} else if err = validator_service.ValidateStrategy(stg, &userOp); err != nil {
		return nil, err
	} else {
		strategy = stg
	}

	//base Strategy and UserOp computeGas
	gasResponse, gasComputeError := gas_service.ComputeGas(&userOp, strategy)
	if gasComputeError != nil {
		return nil, gasComputeError
	}

	//validate gas
	if err := gas_service.ValidateGas(&userOp, gasResponse); err != nil {
		return nil, err
	}
	//pay
	payReceipt, payError := executePay(strategy, &userOp, gasResponse)
	if payError != nil {
		return nil, payError
	}
	paymasterSignature := getPayMasterSignature(strategy, &userOp)
	var result = model.TryPayUserOpResponse{
		StrategyId:         strategy.Id,
		EntryPointAddress:  strategy.EntryPointAddress,
		PayMasterAddress:   strategy.EntryPointAddress,
		PayReceipt:         payReceipt,
		PayMasterSignature: paymasterSignature,
		GasInfo:            gasResponse,
	}

	return &model.Result{
		Code:    200,
		Data:    result,
		Message: "message",
		Cost:    "cost",
	}, nil
}
func businessParamValidate(request *model.TryPayUserOpRequest) error {
	//UserOp Validate
	return nil
}

func executePay(strategy *model.Strategy, userOp *model.UserOperationItem, gasResponse *model.ComputeGasResponse) (interface{}, error) {
	//1.Recharge
	//2.record account
	//3.return Receipt
	return nil, nil
}
func getPayMasterSignature(strategy *model.Strategy, userOp *model.UserOperationItem) string {
	return ""
}

func strategyGenerate(request *model.TryPayUserOpRequest) (*model.Strategy, error) {
	if forceStrategyId := request.ForceStrategyId; forceStrategyId != "" {
		//force strategy
		if strategy := dashboard_service.GetStrategyById(forceStrategyId); strategy == nil {
			return &model.Strategy{}, xerrors.Errorf("Not Support Strategy ID: [%w]", forceStrategyId)
		} else {
			return strategy, nil
		}
	}

	suitableStrategy, err := dashboard_service.GetSuitableStrategy(request.ForceEntryPointAddress, request.ForceNetWork, request.ForceTokens) //TODO
	if err != nil {
		return nil, err
	}
	if suitableStrategy == nil {
		return nil, xerrors.Errorf("Empty Strategies")
	}
	return suitableStrategy, nil
}
