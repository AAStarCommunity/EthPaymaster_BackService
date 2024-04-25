package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
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
func testGetGasPrice(t *testing.T, chain global_const.Network) {
	gasprice, _ := GetGasPrice(chain)
	t.Logf("gasprice %d\n", gasprice.MaxFeePerGas.Uint64())
}

func TestGetAddressTokenBalance(t *testing.T) {
	res, err := GetAddressTokenBalance(global_const.EthereumSepolia, common.HexToAddress("0xDf7093eF81fa23415bb703A685c6331584D30177"), global_const.USDC)
	assert.NoError(t, err)
	fmt.Println(res)
}

func TestChainService(t *testing.T) {
	conf.BasicStrategyInit("../../conf/basic_strategy_dev_config.json")
	conf.BusinessConfigInit("../../conf/business_dev_config.json")
	logrus.SetLevel(logrus.DebugLevel)
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"TestEthereumSepoliaGetPrice",
			func(t *testing.T) {
				testGetGasPrice(t, global_const.EthereumMainnet)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
