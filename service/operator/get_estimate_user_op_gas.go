package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"AAStarCommunity/EthPaymaster_BackService/service/gas_service"
)

func GetEstimateUserOpGas(request *model.UserOpRequest) (*model.ComputeGasResponse, error) {
	var strategy *model.Strategy
	strategy, generateErr := StrategyGenerate(request)
	if generateErr != nil {
		return nil, generateErr
	}

	userOp, err := userop.NewUserOp(&request.UserOp, strategy.GetStrategyEntryPointVersion())
	if err != nil {
		return nil, err
	}
	gasResponse, _, gasComputeError := gas_service.ComputeGas(userOp, strategy)
	if gasComputeError != nil {
		return nil, gasComputeError
	}
	return gasResponse, nil
}
