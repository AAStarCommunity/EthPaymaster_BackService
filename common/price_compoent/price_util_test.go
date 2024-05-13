package price_compoent

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"fmt"
	"strconv"
	"testing"
)

func TestPriceUtilTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.InitConfig("../../config/basic_strategy_config.json", "../../config/basic_config.json", "../../config/secret_config.json")
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"test_ETH_GetPriceUsd",
			func(t *testing.T) {
				testGetPriceUsd(t, global_const.TokenTypeETH)
			},
		},
		{
			"test_OP_GetPriceUsd",
			func(t *testing.T) {
				testGetPriceUsd(t, global_const.TokenTypeOP)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
func testGetPriceUsd(t *testing.T, tokenType global_const.TokenType) {
	price, err := GetPriceUsd(tokenType)
	if err != nil {
		t.Fatal(err)

	}
	t.Logf("price:%v", price)
}

func TestDemo(t *testing.T) {
	str := "0000000000000000000000000000000000000000000000000000000000000002"
	fmt.Printf(strconv.Itoa(len(str)))
}

func TestGetCoinMarketPrice(t *testing.T) {
	GetCoinMarketPrice()
}
