package model

type UserOperationItem struct {
	Sender               string `json:"sender"`
	Nonce                string `json:"nonce"`
	InitCode             string `json:"init_code"`
	CallGasLimit         string `json:"call_gas_limit"`
	VerificationGasList  string `json:"verification_gas_list"`
	PerVerificationGas   string `json:"per_verification_gas"`
	MaxFeePerGas         string `json:"max_fee_per_gas"`
	MaxPriorityFeePerGas string `json:"max_priority_fee_per_gas"`
	//paymasterAndData     string `json:"paymaster_and_data"`
	Signature string `json:"signature"`
}
