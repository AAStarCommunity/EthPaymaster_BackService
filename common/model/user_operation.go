package model

type UserOperationItem struct {
	Sender               string `json:"sender" binding:"required"`
	Nonce                string `json:"nonce" binding:"required"`
	InitCode             string `json:"init_code"`
	CallGasLimit         string `json:"call_gas_limit" binding:"required"`
	VerificationGasList  string `json:"verification_gas_list" binding:"required"`
	PerVerificationGas   string `json:"per_verification_gas" binding:"required"`
	MaxFeePerGas         string `json:"max_fee_per_gas" binding:"required"`
	MaxPriorityFeePerGas string `json:"max_priority_fee_per_gas" binding:"required"`
	Signature            string `json:"signature"`
	//paymasterAndData     string `json:"paymaster_and_data"`
}
