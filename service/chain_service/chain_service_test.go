package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/data_utils"
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/paymaster_data"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"testing"
)

func TestChainService(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.InitConfig("../../config/basic_strategy_config.json", "../../config/basic_config.json", "../../config/secret_config.json")
	logrus.SetLevel(logrus.DebugLevel)
	op, err := user_op.NewUserOp(utils.GenerateMockUservOperation())
	if err != nil {
		t.Error(err)
		return
	}
	opFor1559NotSupport, err := user_op.NewUserOp(utils.GenerateMockUservOperation())
	if err != nil {
		t.Error(err)
		return
	}
	opFor1559NotSupport.MaxPriorityFeePerGas = opFor1559NotSupport.MaxFeePerGas

	tests := []struct {
		name string
		test func(t *testing.T)
	}{

		{
			"TestSepoliaSimulateHandleOp",
			func(t *testing.T) {
				strategy := config.GetBasicStrategyConfig("Ethereum_Sepolia_v06_verifyPaymaster")
				testSimulateHandleOp(t, op, strategy, model.MockGasPrice)
			},
		},
		{
			"TestScrollSepoliaSimulateHandleOp",
			func(t *testing.T) {
				strategy := config.GetBasicStrategyConfig(global_const.StrategyCodeScrollSepoliaV06Verify)
				testSimulateHandleOp(t, opFor1559NotSupport, strategy, model.MockGasPriceNotSupport1559)
			},
		},
		//{ TODO support v07 later
		//	"TestV07SepoliaSimulateHandleOp",
		//	func(t *testing.T) {
		//		strategy := config.GetBasicStrategyConfig(global_const.StrategyCodeEthereumSepoliaV07Verify)
		//		testSimulateHandleOp(t, op, strategy, model.MockGasPrice)
		//	},
		//},
		{
			"testGetpaymasterEntryPointBalance",
			func(t *testing.T) {
				strategy := config.GetBasicStrategyConfig("Ethereum_Sepolia_v06_verifyPaymaster")
				testGetPaymasterEntryPointBalance(t, *strategy)
			},
		},
		{
			"testCheckContractAddressAccess",
			func(t *testing.T) {
				testCheckContractAddressAccess(t)
			},
		},
		{

			"testGetAddressTokenBalance",
			func(t *testing.T) {
				testGetAddressTokenBalance(t)
			},
		},
	}
	for _, tt := range tests {
		if os.Getenv("GITHUB_ACTIONS") != "" && global_const.GitHubActionWhiteListSet.Contains(tt.name) {
			t.Logf("Skip test [%s] in GitHub Actions", tt.name)
			continue
		}
		t.Run(tt.name, tt.test)
	}
}
func testGetAddressTokenBalance(t *testing.T) {
	res, err := GetAddressTokenBalance(global_const.EthereumSepolia, common.HexToAddress("0xFD44DF0Fe211d5EFDBe1423483Fcb3FDeF84540f"), global_const.TokenTypeUSDC)
	assert.NoError(t, err)
	fmt.Println(res)
}
func testCheckContractAddressAccess(t *testing.T) {
	addressStr := "0xF2147CA7f18e8014b76e1A98BaffC96ebB90a29f"
	address := common.HexToAddress(addressStr)
	res, err := CheckContractAddressAccess(&address, global_const.EthereumSepolia)
	assert.NoError(t, err)
	assert.True(t, res)
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

func testSimulateHandleOp(t *testing.T, userOp *user_op.UserOpInput, strategy *model.Strategy, gasPrice *model.GasPrice) {
	paymasterDataInput := paymaster_data.NewPaymasterDataInput(strategy)

	userOpInputForSimulate, err := data_utils.GetUserOpWithPaymasterAndDataForSimulate(*userOp, strategy, paymasterDataInput, gasPrice)
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
