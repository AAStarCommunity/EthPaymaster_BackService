package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckContractAddressAccess(t *testing.T) {
	addressStr := "0x0576a174D229E3cFA37253523E645A78A0C91B57"
	address := common.HexToAddress(addressStr)
	res, err := CheckContractAddressAccess(&address, global_const.EthereumSepolia)
	assert.NoError(t, err)
	assert.True(t, res)
}
func TestGetGasPrice(t *testing.T) {
	gasprice, _ := GetGasPrice(global_const.EthereumMainnet)
	fmt.Printf("gasprice %d\n", gasprice.MaxFeePerGas.Uint64())

}

//	func TestGethClient(t *testing.T) {
//		client, _ := EthCompatibleNetWorkClientMap[global_const.Sepolia]
//		num, _ := client.BlockNumber(context.Background())
//		assert.NotEqual(t, 0, num)
//		fmt.Println(num)
//	}
func TestGetAddressTokenBalance(t *testing.T) {
	res, err := GetAddressTokenBalance(global_const.EthereumSepolia, common.HexToAddress("0xDf7093eF81fa23415bb703A685c6331584D30177"), global_const.USDC)
	assert.NoError(t, err)
	fmt.Println(res)
}
