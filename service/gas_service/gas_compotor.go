package gas_service

import "AAStarCommunity/EthPaymaster_BackService/common/model"

func ComputeGas(userOp *model.UserOperationItem, strategy *model.Strategy) (*model.ComputeGasResponse, error) {
	return &model.ComputeGasResponse{}, nil
}

func ValidateGas(userOp *model.UserOperationItem, gasComputeResponse *model.ComputeGasResponse) error {
	return nil
}
