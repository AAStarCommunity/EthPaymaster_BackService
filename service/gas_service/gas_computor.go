package gas_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
)

func ComputeGas(userOp *model.UserOperationItem, strategy *model.Strategy) (*model.ComputeGasResponse, error) {
	priceInWei, gasPriceInGwei, gasPriceInEtherStr, getGasErr := chain_service.GetGasPrice(types.Sepolia)
	if getGasErr != nil {
		return nil, getGasErr
	}
	return &model.ComputeGasResponse{
		GasPriceInWei:   priceInWei.Uint64(),
		GasPriceInGwei:  gasPriceInGwei,
		GasPriceInEther: *gasPriceInEtherStr,
		TokenCost:       "0.0001",
		Network:         strategy.NetWork,
		Token:           strategy.Token,
		UsdCost:         "0.4",
		BlobEnable:      strategy.Enable4844,
	}, nil
}

func ValidateGas(userOp *model.UserOperationItem, gasComputeResponse *model.ComputeGasResponse) error {
	return nil
}
