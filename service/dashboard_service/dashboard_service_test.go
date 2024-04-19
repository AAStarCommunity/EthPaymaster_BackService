package dashboard_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"fmt"
	"testing"
)

func TestGetSuitableStrategy(t *testing.T) {
	x := types.Network("Ethereum")
	fmt.Println(x)
}
