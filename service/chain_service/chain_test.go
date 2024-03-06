package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckContractAddressAccess(t *testing.T) {
	res, err := CheckContractAddressAccess("0x0576a174D229E3cFA37253523E645A78A0C91B57", types.Sepolia)
	assert.NoError(t, err)
	assert.True(t, res)
}
func TestGetGasPrice(t *testing.T) {
	priceWei, gasPriceInGwei, gasPriceInEtherStr, _ := GetGasPrice(types.Ethereum)
	priceWeiInt := priceWei.Uint64()
	fmt.Printf("priceWeiInt %d\n", priceWeiInt)
	fmt.Printf("gasPriceInGwei %f\n", gasPriceInGwei)
	fmt.Printf("gasPriceInEtherStr %s\n", *gasPriceInEtherStr)

}

func TestGethClient(t *testing.T) {
	client, _ := NetWorkClientMap[types.Sepolia]
	num, _ := client.BlockNumber(context.Background())
	assert.NotEqual(t, 0, num)
	fmt.Println(num)
}