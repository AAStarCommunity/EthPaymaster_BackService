package dashboard_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"fmt"
	"testing"
)

func TestGetSuitableStrategy(t *testing.T) {
	x := global_const.Network("Ethereum")
	fmt.Println(x)
}
