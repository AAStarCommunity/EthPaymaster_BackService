package model

import "math/big"

type GasPrice struct {
	MaxBasePriceWei       *big.Int   `json:"max_base_price_wei"`
	MaxBasePriceGwei      float64    `json:"max_base_price_gwei"`
	MaxBasePriceEther     *big.Float `json:"max_base_price_ether"`
	MaxPriorityPriceWei   *big.Int   `json:"max_priority_price_wei"`
	MaxPriorityPriceGwei  float64    `json:"max_priority_price_gwei"`
	MaxPriorityPriceEther *big.Float `json:"max_priority_price_ether"`
}
