package types

type TokenType string

var StableCoinMap map[TokenType]bool

func init() {
	StableCoinMap = map[TokenType]bool{
		USDT: true,
		USDC: true,
	}
}
func IsStableToken(token TokenType) bool {
	_, ok := StableCoinMap[token]
	return ok
}

const (
	USDT TokenType = "usdt"
	USDC TokenType = "usdc"
	ETH  TokenType = "eth"
	OP   TokenType = "op"
)
