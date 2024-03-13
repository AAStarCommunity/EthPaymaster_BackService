package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckContractAddressAccess(t *testing.T) {
	address := "0x0576a174D229E3cFA37253523E645A78A0C91B57"
	res, err := CheckContractAddressAccess(common.HexToAddress(address), types.Sepolia)
	assert.NoError(t, err)
	assert.True(t, res)
}
func TestGetGasPrice(t *testing.T) {
	gasprice, _ := GetGasPrice(types.Ethereum)
	fmt.Printf("gasprice %d\n", gasprice.MaxBasePriceWei.Uint64())

}

func TestGethClient(t *testing.T) {
	client, _ := EthCompatibleNetWorkClientMap[types.Sepolia]
	num, _ := client.BlockNumber(context.Background())
	assert.NotEqual(t, 0, num)
	fmt.Println(num)
}
