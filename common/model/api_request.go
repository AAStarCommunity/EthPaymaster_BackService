package model

import "errors"

type TryPayUserOpRequest struct {
	ForceStrategyId        string            `json:"strategy_id"`
	ForceNetWork           string            `json:"force_network"`
	ForceTokens            string            `json:"force_tokens"`
	ForceEntryPointAddress string            `json:"force_entry_point_address"`
	UserOperation          UserOperationItem `json:"user_operation"`
	Apikey                 string            `json:"apikey"`
	Extra                  interface{}       `json:"extra"`
}

func (sender *TryPayUserOpRequest) Validate() error {
	if len(sender.Apikey) == 0 {
		return errors.New("apikey mustn't empty")
	}

	if len(sender.ForceStrategyId) == 0 {
		if len(sender.ForceNetWork) == 0 || len(sender.ForceTokens) == 0 || len(sender.ForceEntryPointAddress) == 0 {
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
