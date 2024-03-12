package model

import "math/big"

type GasPrice struct {
	MaxBasePriceWei       *big.Int   `json:"max_base_price_wei"`
	MaxBasePriceGwei      *big.Float `json:"max_base_price_gwei"`
	MaxBasePriceEther     *string    `json:"max_base_price_ether"`
	MaxPriorityPriceWei   *big.Int   `json:"max_priority_price_wei"`
	MaxPriorityPriceGwei  *big.Float `json:"max_priority_price_gwei"`
	MaxPriorityPriceEther *string    `json:"max_priority_price_ether"`
}
