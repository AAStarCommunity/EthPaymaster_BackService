package gas_executor

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/paymaster_data"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"math/big"
	"testing"
)

var (
	MockEstimateGas = &model.UserOpEstimateGas{
		PreVerificationGas:   big.NewInt(52456),
		BaseFee:              big.NewInt(9320437485),
		VerificationGasLimit: big.NewInt(483804),
		CallGasLimit:         big.NewInt(374945),
		MaxFeePerGas:         big.NewInt(10320437485),
		MaxPriorityFeePerGas: big.NewInt(1000000000),
	}
)

func testGetGasPrice(t *testing.T, chain global_const.Network) {
	gasprice, err := GetGasPrice(chain)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("gasprice:%v", gasprice)
}

func TestComputeGas(t *testing.T) {
	//userOp, newErr := user_op.NewUserOp(utils.GenerateMockUservOperation(), global_const.EntrypointV06)
	//assert.NoError(t, newErr)
	//strategy := dashboard_service.GetStrategyById("1")
	//gas, _, err := ComputeGas(userOp, strategy)
	//assert.NoError(t, err)
	//assert.NotNil(t, gas)
	//jsonBypte, _ := json.Marshal(gas)
	//fmt.Println(string(jsonBypte))
	conf.BasicStrategyInit("../conf/basic_strategy_dev_config.json")
	conf.BusinessConfigInit("../conf/business_dev_config.json")
	logrus.SetLevel(logrus.DebugLevel)
	op, err := user_op.NewUserOp(utils.GenerateMockUservOperation())
	if err != nil {
		t.Error(err)
		return
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
			"testEstimateVerificationGasLimit",
			func(*testing.T) {
				strategy := conf.GetBasicStrategyConfig("Ethereum_Sepolia_v06_verifyPaymaster")
				testGetUserOpEstimateGas(t, op, strategy)
			},
		},
		{
			"testEstimateVerificationGasLimit",
			func(*testing.T) {
				totalGasDetail := GetTotalCostByEstimateGas(MockEstimateGas)
				t.Logf("totalGasDetail: %v", totalGasDetail)
				jsonRes, _ := json.Marshal(totalGasDetail)
				t.Logf("totalGasDetail: %v", string(jsonRes))
			},
		},
		{
			"TestGetPreVerificationGas",
			func(t *testing.T) {
				strategy := conf.GetBasicStrategyConfig("Optimism_Sepolia_v06_verifyPaymaster")
				testGetPreVerificationGas(t, op, strategy, model.MockGasPrice)
			},
		},
		{
			"testComputeGas_StrategyCodeEthereumSepoliaVo6Verify",
			func(*testing.T) {
				testComputeGas(t, op, conf.GetBasicStrategyConfig(global_const.StrategyCodeEthereumSepoliaV06Verify))
			},
		},
		{
			"testComputeGas_StrategyCodeOpSepoliaVo6Verify",
			func(*testing.T) {
				testComputeGas(t, op, conf.GetBasicStrategyConfig(global_const.StrategyCodeOptimismSepoliaV06Verify))
			},
		},
		{
			"testComputeGas_StrategyCodeOpSepoliaVo6Verify",
			func(*testing.T) {
				testComputeGas(t, op, conf.GetBasicStrategyConfig(global_const.StrategyCodeOptimismSepoliaV06Verify))
			},
		},
		{
			"testComputeGas_StrategyCodeArbSepoliaVo6Verify",
			func(*testing.T) {
				testComputeGas(t, op, conf.GetBasicStrategyConfig(global_const.StrategyCodeArbitrumSpeoliaV06Verify))
			},
		},
		{
			"testComputeGas_StrategyCodeScrollSepoliaVo6Verify",
			func(*testing.T) {
				testComputeGas(t, op, conf.GetBasicStrategyConfig(global_const.StrategyCodeScrollSepoliaV06Verify))
			},
		},
		{
			"testComputeGas_StrategyCodeBaseSepoliaVo6Verify",
			func(*testing.T) {
				testComputeGas(t, op, conf.GetBasicStrategyConfig(global_const.StrategyCodeBaseSepoliaV06Verify))
			},
		},
		{
			"testComputeGas_StrategyCodeEthereumSepoliaV06Erc20",
			func(*testing.T) {
				strategy := conf.GetBasicStrategyConfig(global_const.StrategyCodeEthereumSepoliaV06Erc20)
				strategy.Erc20TokenType = global_const.TokenTypeUSDT
				testComputeGas(t, op, strategy)
			},
		},
		{
			"testComputeGas_StrategyCodeOptimismSepoliaV06Erc20",
			func(*testing.T) {
				strategy := conf.GetBasicStrategyConfig(global_const.StrategyCodeOptimismSepoliaV06Erc20)
				strategy.Erc20TokenType = global_const.TokenTypeUSDT
				testComputeGas(t, op, strategy)
			},
		},
		{
			"testComputeGas_StrategyCodeArbitrumSpeoliaV06Erc20",
			func(*testing.T) {
				strategy := conf.GetBasicStrategyConfig(global_const.StrategyCodeArbitrumSpeoliaV06Erc20)
				strategy.Erc20TokenType = global_const.TokenTypeUSDT
				testComputeGas(t, op, strategy)
			},
		},
		{
			"testComputeGas_StrategyCodeBaseSepoliaV06Erc20",
			func(*testing.T) {
				strategy := conf.GetBasicStrategyConfig(global_const.StrategyCodeBaseSepoliaV06Erc20)
				strategy.Erc20TokenType = global_const.TokenTypeUSDT
				testComputeGas(t, op, strategy)
			},
		},
		{
			"TestScrollEstimateCallGasLimit",
			func(t *testing.T) {
				testEstimateCallGasLimit(t, conf.GetBasicStrategyConfig(global_const.StrategyCodeScrollSepoliaV06Verify), model.MockSimulateHandleOpResult, op, global_const.DummayPreverificationgasBigint)
			},
		},
		{
			"TestErc20",
			func(t *testing.T) {
				strategy := conf.GetBasicStrategyConfig(global_const.StrategyCodeEthereumSepoliaV06Erc20)
				strategy.Erc20TokenType = global_const.TokenTypeUSDT
				testErc20TokenCost(t, strategy, big.NewFloat(0.0001))
			},
		},
		{
			"TestUSDTTokenCost",
			func(t *testing.T) {
				strategy := conf.GetBasicStrategyConfig(global_const.StrategyCodeEthereumSepoliaV06Erc20)
				strategy.Erc20TokenType = global_const.TokenTypeUSDT
				testErc20TokenCost(t, strategy, big.NewFloat(0.0001))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
func testErc20TokenCost(t *testing.T, strategy *model.Strategy, tokenCount *big.Float) {
	erc20TokenCost, err := getErc20TokenCost(strategy, tokenCount)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("erc20TokenCost:%v", erc20TokenCost)

}
func testEstimateCallGasLimit(t *testing.T, strategy *model.Strategy, simulateOpResult *model.SimulateHandleOpResult, op *user_op.UserOpInput, simulateGasPrice *big.Int) {
	callGasLimit, err := EstimateCallGasLimit(strategy, simulateOpResult, op, simulateGasPrice)
	if err != nil {
		t.Error(err)
		return
	}
	if callGasLimit == nil {
		t.Error("callGasLimit is nil")
		return
	}
	t.Logf("callGasLimit: %v", callGasLimit)

}
func testGetPreVerificationGas(t *testing.T, userOp *user_op.UserOpInput, strategy *model.Strategy, gasFeeResult *model.GasPrice) {
	res, err := GetPreVerificationGas(userOp, strategy, gasFeeResult, nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("preVerificationGas:%v", res)
}
func testComputeGas(t *testing.T, input *user_op.UserOpInput, strategy *model.Strategy) {
	t.Logf("strategy: %v", strategy)
	paymasterDataInput := paymaster_data.NewPaymasterDataInput(strategy)
	res, _, err := ComputeGas(input, strategy, paymasterDataInput)
	if err != nil {
		t.Error(err)
		return
	}
	if res == nil {
		t.Error("res is nil")
		return
	}
	jsonRes, _ := json.Marshal(res)
	t.Logf("res: %v", string(jsonRes))
}
func TestEstimateCallGasLimit(t *testing.T) {
	callGasLimit, err := estimateVerificationGasLimit(model.MockSimulateHandleOpResult, global_const.DummayPreverificationgasBigint)

	if err != nil {
		t.Error(err)
		return
	}
	if callGasLimit == nil {
		t.Error("callGasLimit is nil")
		return
	}
	t.Logf("callGasLimit: %v", callGasLimit)
}
func testGetUserOpEstimateGas(t *testing.T, input *user_op.UserOpInput, strategy *model.Strategy) {
	paymasterDataInput := paymaster_data.NewPaymasterDataInput(strategy)
	res, err := getUserOpEstimateGas(input, strategy, paymasterDataInput)
	if err != nil {
		t.Error(err)
		return
	}
	if res == nil {
		t.Error("res is nil")
		return
	}
	jsonRes, _ := json.Marshal(res)
	t.Logf("res: %v", string(jsonRes))
}
func testBase(t *testing.T) {
	client, err := ethclient.Dial(conf.GetEthereumRpcUrl("https://base-sepolia.g.alchemy.com/v2/zUhtd18b2ZOTIJME6rv2Uwz9q7PBnnsa"))
	if err != nil {
		t.Error(err)
		return
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("gasPrice:%v", gasPrice)
}
