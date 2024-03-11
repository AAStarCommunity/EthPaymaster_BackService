package model

type UserOperationItem struct {
	Sender               string `json:"sender" binding:"required,hexParam"`
	Nonce                string `json:"nonce" binding:"required"`
	InitCode             string `json:"init_code"`
	CallData             string `json:"call_data" binding:"required"`
	CallGasLimit         string `json:"call_gas_limit" binding:"required"`
	VerificationGasList  string `json:"verification_gas_list" binding:"required"`
	PreVerificationGas   string `json:"per_verification_gas" binding:"required"`
	MaxFeePerGas         string `json:"max_fee_per_gas" binding:"required"`
	MaxPriorityFeePerGas string `json:"max_priority_fee_per_gas" binding:"required"`
	Signature            string `json:"signature" binding:"required"`
	PaymasterAndData     string `json:"paymaster_and_data"`
}
