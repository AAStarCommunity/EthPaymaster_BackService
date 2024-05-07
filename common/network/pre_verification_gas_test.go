package network

import (
	"AAStarCommunity/EthPaymaster_BackService/config"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestPreVerGas(t *testing.T) {
	config.BasicStrategyInit("../../config/basic_strategy_dev_config.json")
	config.BusinessConfigInit("../../config/business_dev_config.json")
	logrus.SetLevel(logrus.DebugLevel)

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
