package utils

import (
	"AAStarCommunity/EthPaymaster_BackService/common/tokens"
	"fmt"
	"strconv"
	"testing"
)

func TestGetPriceUsd(t *testing.T) {
	price, _ := GetPriceUsd(tokens.OP)
	fmt.Println(price)
}

func TestGetToken(t *testing.T) {
	price, _ := GetToken(tokens.ETH, tokens.OP)
	fmt.Println(price)
}
func TestDemo(t *testing.T) {
	str := "0000000000000000000000000000000000000000000000000000000000000002"
	fmt.Printf(strconv.Itoa(len(str)))
}

func TestGetCoinMarketPrice(t *testing.T) {
	GetCoinMarketPrice()
}
