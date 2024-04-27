package network

import (
	"context"
	"fmt"
	"github.com/NethermindEth/starknet.go/rpc"
	"testing"
)

func TestDemo(t *testing.T) {
	starkProvider, err := rpc.NewProvider("https://starknet-sepolia.infura.io/v3/0284f5a9fc55476698079b24e2f97909")
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
