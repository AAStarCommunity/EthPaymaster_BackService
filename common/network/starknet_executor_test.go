package network

import (
	"context"
	starknetRpc "github.com/NethermindEth/starknet.go/rpc"
	"testing"
)

func TestDemo(t *testing.T) {
	starkProvider, err := starknetRpc.NewProvider("")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	chainId, err := starkProvider.ChainID(context.Background())
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	t.Logf("Chain ID: %v", chainId)
}
