package model

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"errors"
)

type TryPayUserOpRequest struct {
	ForceStrategyId        string            `json:"force_strategy_id"`
	ForceNetwork           types.NetWork     `json:"force_network"`
	ForceToken             string            `json:"force_token"`
	ForceEntryPointAddress string            `json:"force_entrypoint_address"`
	UserOperation          UserOperationItem `json:"user_operation"`
	Extra                  interface{}       `json:"extra"`
}

func (request *TryPayUserOpRequest) Validate() error {
	if len(request.ForceStrategyId) == 0 {
		if len(request.ForceNetwork) == 0 || len(request.ForceToken) == 0 || len(request.ForceEntryPointAddress) == 0 {
			return errors.New("strategy configuration illegal")
		}
	}
	return nil
}

type GetSupportEntrypointRequest struct {
}

func (request *GetSupportEntrypointRequest) Validate() error {
	return nil
}

type GetSupportStrategyRequest struct {
}

func (request *GetSupportStrategyRequest) Validate() error {
	return nil
}
