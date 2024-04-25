package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestOperator(t *testing.T) {
	conf.BasicStrategyInit("../../conf/basic_strategy_dev_config.json")
	conf.BusinessConfigInit("../../conf/business_dev_config.json")
	logrus.SetLevel(logrus.DebugLevel)
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"testGetSupportEntrypointExecute",
			func(t *testing.T) {
				testGetSupportEntrypointExecute(t)
			},
		},
		{
			"TestTryPayUserOpExecute",
			func(t *testing.T) {
				testTryPayUserOpExecute(t)
			},
		},
		{
			"testGetSupportStrategyExecute",
			func(t *testing.T) {
				testGetSupportStrategyExecute(t)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}

}
func testGetSupportEntrypointExecute(t *testing.T) {
	res, err := GetSupportEntrypointExecute("network")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}
func testTryPayUserOpExecute(t *testing.T) {
	request := getMockTryPayUserOpRequest()
	result, err := TryPayUserOpExecute(request)
	if err != nil {
		t.Error(err)
		return
	}
	resultJson, _ := json.Marshal(result)
	fmt.Printf("Result: %v", string(resultJson))
}

func getMockTryPayUserOpRequest() *model.UserOpRequest {
	return &model.UserOpRequest{
		ForceStrategyId: "1",
		UserOp:          *utils.GenerateMockUservOperation(),
	}
}

func testGetSupportStrategyExecute(t *testing.T) {
	res, err := GetSupportStrategyExecute("network")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)

}
