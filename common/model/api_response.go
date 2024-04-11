package model

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"math/big"
)

type TryPayUserOpResponse struct {
	StrategyId         string              `json:"strategy_id"`
	EntryPointAddress  string              `json:"entrypoint_address"`
	PayMasterAddress   string              `json:"paymaster_address"`
	PayMasterSignature string              `json:"paymaster_signature"`
	PayMasterAndData   string              `json:"paymaster_and_data"`
	PayReceipt         *PayReceipt         `json:"pay_receipt"`
	GasInfo            *ComputeGasResponse `json:"gas_info"`
}

type ComputeGasResponse struct {
	GasInfo       *GasPrice          `json:"gas_info"`
	TokenCost     *big.Float         `json:"token_cost"`
	Network       types.Network      `json:"network"`
	Token         types.TokenType    `json:"tokens"`
	UsdCost       float64            `json:"usd_cost"`
	BlobEnable    bool               `json:"blob_enable"`
	MaxFee        big.Int            `json:"max_fee"`
	OpEstimateGas *UserOpEstimateGas `json:"op_estimate_gas"`
}
type UserOpEstimateGas struct {
	PreVerificationGas   *big.Int `json:"preVerificationGas"`
	VerificationGasLimit *big.Int `json:"verificationGasLimit"`
	CallGasLimit         *big.Int `json:"callGasLimit"`
	VerificationGas      *big.Int `json:"verificationGas"`
	MaxFeePerGas         *big.Int `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *big.Int `json:"maxPriorityFeePerGas"`
	//v0.7
	AccountGasLimit               string `json:"account_gas_limit" binding:"required"`
	PaymasterVerificationGasLimit string `json:"paymaster_verification_gas_limit" binding:"required"`
	PaymasterPostOpGasLimit       string `json:"paymaster_post_op_gas_limit" binding:"required"`
	GasFees                       []byte `json:"gasFees" binding:"required"`
}
type PayReceipt struct {
	TransactionHash string `json:"transaction_hash"`
	Sponsor         string `json:"sponsor"`
}

type GetSupportEntryPointResponse struct {
	EntrypointDomains *[]EntrypointDomain `json:"entrypoints"`
}
type EntrypointDomain struct {
	Address    string        `json:"address"`
	Desc       string        `json:"desc"`
	NetWork    types.Network `json:"network"`
	StrategyId string        `json:"strategy_id"`
}

type GetSupportStrategyResponse struct {
	Strategies *[]Strategy `json:"strategies"`
}
