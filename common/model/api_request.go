package model

import "errors"

type TryPayUserOpRequest struct {
	ForceStrategyId        string            `json:"strategy_id"`
	ForceNetWork           string            `json:"force_network"`
	ForceTokens            string            `json:"force_tokens"`
	ForceEntryPointAddress string            `json:"force_entry_point_address"`
	UserOperation          UserOperationItem `json:"user_operation"`
	Extra                  interface{}       `json:"extra"`
}

func (request *TryPayUserOpRequest) Validate() error {
	if len(request.ForceStrategyId) == 0 {
		if len(request.ForceNetWork) == 0 || len(request.ForceTokens) == 0 || len(request.ForceEntryPointAddress) == 0 {
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
