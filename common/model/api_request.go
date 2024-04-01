package model

import (
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"errors"
)

type TryPayUserOpRequest struct {
	ForceStrategyId        string          `json:"force_strategy_id"`
	ForceNetwork           network.Network `json:"force_network"`
	ForceToken             string          `json:"force_token"`
	ForceEntryPointAddress string          `json:"force_entrypoint_address"`
	UserOp                 map[string]any  `json:"user_operation"`
	Extra                  interface{}     `json:"extra"`
	OnlyEstimateGas        bool            `json:"only_estimate_gas"`
}

func (request *TryPayUserOpRequest) Validate() error {
	if len(request.ForceStrategyId) == 0 {
		if len(request.ForceNetwork) == 0 || len(request.ForceToken) == 0 || len(request.ForceEntryPointAddress) == 0 {
			return errors.New("strategy configuration illegal")
		}
	}
	return nil
}
