package network

import (
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestPreVerGas(t *testing.T) {
	conf.BasicStrategyInit("../../conf/basic_strategy_dev_config.json")
	conf.BusinessConfigInit("../../conf/business_dev_config.json")
	logrus.SetLevel(logrus.DebugLevel)
	//op, err := user_op.NewUserOp(utils.GenerateMockUservOperation())
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"TestOpPreVerGas",
			func(t *testing.T) {

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

//func testL2PreVerGas(t *testing.T,) {
//
//}
