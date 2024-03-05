package model

type TryPayUserOpRequest struct {
	ForceStrategyId        string            `json:"strategy_id"`
	ForceNetWork           string            `json:"force_network"`
	ForceTokens            string            `json:"force_tokens"`
	ForceEntryPointAddress string            `json:"force_entry_point_address"`
	UserOperation          UserOperationItem `json:"user_operation"`
	Apikey                 string            `json:"apikey"`
	Extra                  interface{}       `json:"extra"`
}
