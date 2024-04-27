package gas_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/paymaster_data"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"encoding/json"
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

func TestComputeGas(t *testing.T) {
	//userOp, newErr := user_op.NewUserOp(utils.GenerateMockUservOperation(), global_const.EntrypointV06)
	//assert.NoError(t, newErr)
	//strategy := dashboard_service.GetStrategyById("1")
	//gas, _, err := ComputeGas(userOp, strategy)
	//assert.NoError(t, err)
	//assert.NotNil(t, gas)
	//jsonBypte, _ := json.Marshal(gas)
	//fmt.Println(string(jsonBypte))
	conf.BasicStrategyInit("../../conf/basic_strategy_dev_config.json")
	conf.BusinessConfigInit("../../conf/business_dev_config.json")
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
			"testComputeGas",
			func(*testing.T) {
				testComputeGas(t, op, conf.GetBasicStrategyConfig("Ethereum_Sepolia_v06_verifyPaymaster"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}

}
func testComputeGas(t *testing.T, input *user_op.UserOpInput, strategy *model.Strategy) {
	paymasterDataInput := paymaster_data.NewPaymasterDataInput(strategy)
	res, _, err := ComputeGas(input, strategy, paymasterDataInput)
	if err != nil {
		logrus.Error(err)
		return
	}
	if res == nil {
		logrus.Error("res is nil")
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
