package global_const

import (
	mapset "github.com/deckarep/golang-set/v2"
)

type TokenType string

var StableCoinSet mapset.Set[TokenType]

//var StableCoinMap map[TokenType]bool

func init() {
	StableCoinSet = mapset.NewSet[TokenType]()
	StableCoinSet.Add(USDT)
	StableCoinSet.Add(USDC)

}
func IsStableToken(token TokenType) bool {
	return StableCoinSet.Contains(token)

}

const (
	USDT TokenType = "USDT"
	USDC TokenType = "USDC"
	ETH  TokenType = "ETH"
	OP   TokenType = "OP"
)
