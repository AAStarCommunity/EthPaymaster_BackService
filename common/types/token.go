package types

import (
	mapset "github.com/deckarep/golang-set/v2"
)

type TokenType string

var StableCoinSet mapset.Set[TokenType]

//var StableCoinMap map[TokenType]bool

func init() {
	StableCoinSet.Add(USDT)
	StableCoinSet.Add(USDC)

}
func IsStableToken(token TokenType) bool {
	return StableCoinSet.Contains(token)

}

const (
	USDT TokenType = "usdt"
	USDC TokenType = "usdc"
	ETH  TokenType = "eth"
	OP   TokenType = "op"
)
