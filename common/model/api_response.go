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
	GasInfo    *GasPrice       `json:"gas_info"`
	TokenCost  *big.Float      `json:"token_cost"`
	Network    types.Network   `json:"network"`
	Token      types.TokenType `json:"token"`
	UsdCost    float64         `json:"usd_cost"`
	BlobEnable bool            `json:"blob_enable"`
	MaxFee     big.Int         `json:"max_fee"`
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
