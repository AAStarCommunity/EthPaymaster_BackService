package model

import "math/big"

var (
	MockSimulateHandleOpResult = &SimulateHandleOpResult{
		PreOpGas:      big.NewInt(374992),
		GasPaid:       big.NewInt(595492500000000),
		ValidAfter:    big.NewInt(1710044496),
		ValidUntil:    big.NewInt(1820044496),
		TargetSuccess: false,
		TargetResult:  nil,
	}
)

type GasPrice struct {
	MaxFeePerGas *big.Int `json:"max_base_price_wei"`
	//MaxBasePriceGwei      float64    `json:"max_base_price_gwei"`
	//MaxBasePriceEther     *big.Float `json:"max_base_price_ether"`
	MaxPriorityFeePerGas *big.Int `json:"max_priority_price_wei"`
	//MaxPriorityPriceGwei  float64    `json:"max_priority_price_gwei"`
	//MaxPriorityPriceEther *big.Float `json:"max_priority_price_ether"`
	BaseFee *big.Int `json:"base_fee"`
}
type SimulateHandleOpResult struct {
	// PreOpGas = preGas - gasleft() + userOp.preVerificationGas;
	// PreOpGas = verificationGasLimit + userOp.preVerificationGas;
	PreOpGas      *big.Int `json:"preOpGas"`
	GasPaid       *big.Int `json:"paid"`
	ValidAfter    *big.Int `json:"validAfter"`
	ValidUntil    *big.Int `json:"validUntil"`
	TargetSuccess bool
	TargetResult  []byte
}

type GasFeePerGasResult struct {
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
	BaseFee              *big.Int
}
