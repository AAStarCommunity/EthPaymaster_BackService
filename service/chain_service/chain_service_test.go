package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/data_utils"
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/paymaster_data"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"math/big"
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
	gasprice, err := GetGasPrice(chain)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("gasprice:%v", gasprice)
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
	op, err := user_op.NewUserOp(utils.GenerateMockUservOperation())
	if err != nil {
		t.Error(err)
		return
	}
	mockGasPrice := &model.GasPrice{
		MaxFeePerGas:         big.NewInt(2053608903),
		MaxPriorityFeePerGas: big.NewInt(1000000000),
		BaseFee:              big.NewInt(1053608903),
	}
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"TestEthereumSepoliaGetPrice",
			func(t *testing.T) {
				testGetGasPrice(t, global_const.EthereumSepolia)
			},
		},
		{
			"TestGetPreVerificationGas",
			func(t *testing.T) {
				strategy := conf.GetBasicStrategyConfig("Optimism_Sepolia_v06_verifyPaymaster")
				testGetPreVerificationGas(t, op, strategy, mockGasPrice)
			},
		},
		{
			"TestSepoliaSimulateHandleOp",
			func(t *testing.T) {
				strategy := conf.GetBasicStrategyConfig("Ethereum_Sepolia_v06_verifyPaymaster")
				testSimulateHandleOp(t, op, strategy)
			},
		},
		{
			"testGetpaymasterEntryPointBalance",
			func(t *testing.T) {
				strategy := conf.GetBasicStrategyConfig("Ethereum_Sepolia_v06_verifyPaymaster")
				testGetPaymasterEntryPointBalance(t, *strategy)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
func testGetPaymasterEntryPointBalance(t *testing.T, strategy model.Strategy) {
	res, err := GetPaymasterEntryPointBalance(&strategy)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(res)
	t.Logf("paymasterEntryPointBalance:%v", res)

}
func testGetPreVerificationGas(t *testing.T, userOp *user_op.UserOpInput, strategy *model.Strategy, gasFeeResult *model.GasPrice) {
	res, err := GetPreVerificationGas(userOp, strategy, gasFeeResult)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("preVerificationGas:%v", res)
}
func testSimulateHandleOp(t *testing.T, userOp *user_op.UserOpInput, strategy *model.Strategy) {
	paymasterDataInput := paymaster_data.NewPaymasterDataInput(strategy)
	userOpInputForSimulate, err := data_utils.GetUserOpWithPaymasterAndDataForSimulate(*userOp, strategy, paymasterDataInput, &model.GasPrice{})
	if err != nil {
		t.Error(err)
		return
	}
	res, err := SimulateHandleOp(userOpInputForSimulate, strategy)
	if err != nil {
		t.Error(err)
		return
	}
	jsonRes, _ := json.Marshal(res)
	t.Logf("simulateHandleOp:%v", string(jsonRes))
	callGasCount := new(big.Int).Div(res.GasPaid, res.PreOpGas)
	t.Logf("callGasCount:%v", callGasCount)
}
