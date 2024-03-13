package utils

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"fmt"
	"testing"
)

func TestGetPriceUsd(t *testing.T) {
	price, _ := GetPriceUsd(types.OP)
	fmt.Println(price)
}

func TestGetToken(t *testing.T) {
	price, _ := GetToken(types.ETH, types.OP)
	fmt.Println(price)
}
