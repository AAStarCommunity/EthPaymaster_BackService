package model

type TryPayUserOpResponse struct {
	StrategyId         string              `json:"strategy_id"`
	EntryPointAddress  string              `json:"entry_point_address"`
	PayMasterAddress   string              `json:"pay_master_address"`
	PayMasterSignature string              `json:"pay_master_signature"`
	PayReceipt         interface{}         `json:"pay_receipt"`
	GasInfo            *ComputeGasResponse `json:"gaf_info"`
}

type ComputeGasResponse struct {
	StrategyId string
	TokenCost  string
	Network    string
	Token      string
	UsdCost    string
}
