package utils

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"fmt"
	"strconv"
	"testing"
)

func TestGetPriceUsd(t *testing.T) {
	price, _ := GetPriceUsd(global_const.OP)
	fmt.Println(price)
}

func TestGetToken(t *testing.T) {
	price, _ := GetToken(global_const.ETH, global_const.OP)
	fmt.Println(price)
}
func TestDemo(t *testing.T) {
	str := "0000000000000000000000000000000000000000000000000000000000000002"
	fmt.Printf(strconv.Itoa(len(str)))
}

func TestGetCoinMarketPrice(t *testing.T) {
	GetCoinMarketPrice()
}
