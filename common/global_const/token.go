package global_const

import (
	mapset "github.com/deckarep/golang-set/v2"
)

type TokenType string

var StableCoinSet mapset.Set[TokenType]

//var StableCoinMap map[TokenType]bool

func init() {
	StableCoinSet = mapset.NewSet[TokenType]()
	StableCoinSet.Add(TokenTypeUSDT)
	StableCoinSet.Add(TokenTypeUSDC)

}
func IsStableToken(token TokenType) bool {
	return StableCoinSet.Contains(token)

}

const (
	TokenTypeUSDT   TokenType = "USDT"
	TokenTypeUSDC   TokenType = "USDC"
	TokenTypeETH    TokenType = "ETH"
	TokenTypeOP     TokenType = "OP"
	TokenTypeAAStar TokenType = "AAStar"
)
