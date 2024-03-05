package model

type TryPayUserOpResponse struct {
	StrategyId         string              `json:"strategy_id"`
	EntryPointAddress  string              `json:"entrypoint_address"`
	PayMasterAddress   string              `json:"paymaster_address"`
	PayMasterSignature string              `json:"paymaster_signature"`
	PayReceipt         interface{}         `json:"pay_receipt"`
	GasInfo            *ComputeGasResponse `json:"gas_info"`
}

type ComputeGasResponse struct {
	StrategyId string `json:"strategy_id"`
	TokenCost  string `json:"token_cost"`
	Network    string `json:"network"`
	Token      string `json:"token"`
	UsdCost    string `json:"usd_cost"`
}
