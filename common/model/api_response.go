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
	PayReceipt         *PayReceipt         `json:"pay_receipt"`
	GasInfo            *ComputeGasResponse `json:"gas_info"`
}

type ComputeGasResponse struct {
	GasPriceInWei   uint64          `json:"gas_price_wei"` // wei
	GasPriceInGwei  *big.Float      `json:"gas_price_gwei"`
	GasPriceInEther string          `json:"gas_price_ether"`
	TokenCost       string          `json:"token_cost"`
	Network         types.NetWork   `json:"network"`
	Token           types.TokenType `json:"token"`
	UsdCost         string          `json:"usd_cost"`
	BlobEnable      bool            `json:"blob_enable"`
}
type PayReceipt struct {
	TransactionHash string `json:"transaction_hash"`
	Sponsor         string `json:"sponsor"`
}
