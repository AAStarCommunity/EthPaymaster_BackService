package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckContractAddressAccess(t *testing.T) {
	addressStr := "0x0576a174D229E3cFA37253523E645A78A0C91B57"
	address := common.HexToAddress(addressStr)
	res, err := CheckContractAddressAccess(&address, types.EthereumSepolia)
	assert.NoError(t, err)
	assert.True(t, res)
}
func TestGetGasPrice(t *testing.T) {
	gasprice, _ := GetGasPrice(types.EthereumMainnet)
	fmt.Printf("gasprice %d\n", gasprice.MaxFeePerGas.Uint64())

}

//	func TestGethClient(t *testing.T) {
//		client, _ := EthCompatibleNetWorkClientMap[types.Sepolia]
//		num, _ := client.BlockNumber(context.Background())
//		assert.NotEqual(t, 0, num)
//		fmt.Println(num)
//	}
func TestGetAddressTokenBalance(t *testing.T) {

	res, err := GetAddressTokenBalance(types.EthereumSepolia, common.HexToAddress("0xDf7093eF81fa23415bb703A685c6331584D30177"), types.USDC)
	assert.NoError(t, err)
	fmt.Println(res)
}
