package sponsor_manager

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"math/big"
	"testing"
)

func TestSponsor(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.InitConfig("../config/basic_strategy_config.json", "../config/basic_config.json", "../config/secret_config.json")
	Init()
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"TestGetAvailableBalance",
			func(t *testing.T) {
				request := model.DepositSponsorRequest{
					PayUserId: "test",
					Amount:    big.NewFloat(1),
					IsTestNet: true,
				}
				result, err := DepositSponsor(&request)
				if err != nil {
					t.Error(err)
				}
				if result == nil {
					t.Error("DepositSponsor failed")
				}
			},
		},
		{
			"TestLockUserBalanceSt",
			func(t *testing.T) {
				res, err := SelectUserSponsorBalanceDBModelWithScanList()
				if err != nil {
					t.Error(err)
				}
				if res == nil {
					t.Error("SelectUserSponsorBalanceDBModelWithScanList failed")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
