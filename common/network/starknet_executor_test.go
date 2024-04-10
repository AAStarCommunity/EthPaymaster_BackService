package network

import (
	"context"
	"fmt"
	"github.com/NethermindEth/starknet.go/rpc"
	"testing"
)

func TestDemo(t *testing.T) {
	starkProvider, err := rpc.NewProvider("https://starknet-sepolia.g.alchemy.com/v2/uuXjaVAZy6-uzgoobtYd1IIX-IfjvXBc")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	chainId, chainIdError := starkProvider.ChainID(context.Background())
	if chainIdError != nil {
		t.Errorf("Error: %v", chainIdError)
		return
	}
	fmt.Println(chainId)
}
