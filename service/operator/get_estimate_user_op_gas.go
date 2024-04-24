package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/service/gas_service"
)

func GetEstimateUserOpGas(request *model.UserOpRequest) (*model.ComputeGasResponse, error) {
	var strategy *model.Strategy
	strategy, generateErr := StrategyGenerate(request)
	if generateErr != nil {
		return nil, generateErr
	}

	userOp, err := user_op.NewUserOp(&request.UserOp)
	if err != nil {
		return nil, err
	}
	gasResponse, _, gasComputeError := gas_service.ComputeGas(userOp, strategy, model.NewPaymasterDataInput(strategy))
	if gasComputeError != nil {
		return nil, gasComputeError
	}
	return gasResponse, nil
}
