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

	fmt.Printf("gaspricegwei %f\n", gasprice.MaxBasePriceGwei)
	fmt.Printf("gaspriceeth %s\n", gasprice.MaxBasePriceEther.String())

}

func TestGethClient(t *testing.T) {
	client, _ := EthCompatibleNetWorkClientMap[types.Sepolia]
	num, _ := client.BlockNumber(context.Background())
	assert.NotEqual(t, 0, num)
	fmt.Println(num)
}
func TestGetAddressTokenBalance(t *testing.T) {

	res, err := GetAddressTokenBalance(types.Sepolia, common.HexToAddress("0xDf7093eF81fa23415bb703A685c6331584D30177"), types.USDC)
	assert.NoError(t, err)
	fmt.Println(res)
}

func TestGetChainId(t *testing.T) {
	res, err := GetChainId(types.Sepolia)
	assert.NoError(t, err)
	fmt.Println(res)
}
