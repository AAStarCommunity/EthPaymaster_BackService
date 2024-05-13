package starknet

import (
	"context"
	"fmt"
	"github.com/NethermindEth/starknet.go/rpc"
	"testing"
)

func TestDemo(t *testing.T) {
	//only read
	starkProvider, err := rpc.NewProvider("https://starknet-sepolia.public.blastapi.io/rpc/v0_7")
	if err != nil {
		t.Fatalf("Error: %v", err)
		return
	}
	chainId, chainIdError := starkProvider.ChainID(context.Background())
	if chainIdError != nil {
		t.Fatalf("Error: %v", chainIdError)
		return
	}
	//starkProvider.SimulateTransactions()
	fmt.Println(chainId)

}
