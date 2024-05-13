package model

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"math/big"
)

type TryPayUserOpResponse struct {
	StrategyId        string                         `json:"strategyId"`
	NetWork           global_const.Network           `json:"network"`
	EntrypointVersion global_const.EntrypointVersion `json:"entrypointVersion"`
	EntryPointAddress string                         `json:"entrypointAddress"`
	PayMasterAddress  string                         `json:"paymasterAddress"`
	Erc20TokenCost    *big.Float                     `json:"Erc20TokenCost"`
	UserOpResponse    *UserOpResponse                `json:"userOpResponse"`
}
type UserOpResponse struct {
	PayMasterAndData     string   `json:"paymasterAndData"`
	PreVerificationGas   *big.Int `json:"preVerificationGas"`
	VerificationGasLimit *big.Int `json:"verificationGasLimit"`
	CallGasLimit         *big.Int `json:"callGasLimit"`
	MaxFeePerGas         *big.Int `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *big.Int `json:"maxPriorityFeePerGas"`
	//v0.7
	AccountGasLimit               string   `json:"accountGasLimit" binding:"required"`
	PaymasterVerificationGasLimit *big.Int `json:"paymasterVerificationGasLimit" binding:"required"`
	PaymasterPostOpGasLimit       *big.Int `json:"paymasterPostOpGasLimit" binding:"required"`
	GasFees                       string   `json:"gasFees" binding:"required"`
}

type ComputeGasResponse struct {
	Erc20TokenCost *big.Float         `json:"Erc20TokenCost"`
	OpEstimateGas  *UserOpEstimateGas `json:"opEstimateGas"`
	TotalGasDetail *TotalGasDetail    `json:"totalGasDetail"`
}
type UserOpEstimateGas struct {
	//common
	PreVerificationGas *big.Int `json:"preVerificationGas"`

	BaseFee *big.Int `json:"baseFee"`
	//v0.6
	VerificationGasLimit *big.Int `json:"verificationGasLimit"`
	CallGasLimit         *big.Int `json:"callGasLimit"`
	MaxFeePerGas         *big.Int `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *big.Int `json:"maxPriorityFeePerGas"`
	//v0.7
	AccountGasLimit               *[32]byte `json:"accountGasLimit" binding:"required"`
	PaymasterVerificationGasLimit *big.Int  `json:"paymasterVerificationGasLimit" binding:"required"`
	PaymasterPostOpGasLimit       *big.Int  `json:"paymasterPostOpGasLimit" binding:"required"`
	GasFees                       *[32]byte `json:"gasFees" binding:"required"`
}
type PayReceipt struct {
	TransactionHash string `json:"transaction_hash"`
	Sponsor         string `json:"sponsor"`
}

type GetSupportEntryPointResponse struct {
	EntrypointDomains *[]EntrypointDomain `json:"entrypoints"`
}
type EntrypointDomain struct {
	Address    string               `json:"address"`
	Desc       string               `json:"desc"`
	NetWork    global_const.Network `json:"network"`
	StrategyId string               `json:"strategy_id"`
}

type GetSupportStrategyResponse struct {
	Strategies *[]Strategy `json:"strategies"`
}
