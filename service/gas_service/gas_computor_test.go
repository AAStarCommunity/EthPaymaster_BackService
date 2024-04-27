package gas_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/paymaster_data"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"github.com/sirupsen/logrus"
	"testing"
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
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}

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
	t.Logf("res: %v", res)
}
