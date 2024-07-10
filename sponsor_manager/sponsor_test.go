package sponsor_manager

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/config"

	"encoding/json"
	"math/big"
	"testing"
)

func TestSponsor(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.InitConfig("../config/basic_strategy_config.json", "../config/basic_config.json", "../config/secret_config.json")
	Init()
	mockUserOpHash, err := utils.DecodeStringWithPrefix("0xc0977a4ed200e8c2e62b906260e57ee2f7dca962089f8a4645117bfee3b1f215")
	mockUserOpHash2, err := utils.DecodeStringWithPrefix("0xc0977a4ed200e8c2e62b906260e57ee2f7dca962089f8a4645117bfee3b1f211")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"TestDepositSponsor",
			func(t *testing.T) {
				request := DepositSponsorInput{
					PayUserId: "test",
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
			"TestGetAvaliabeBalance",
			func(t *testing.T) {
				testFindUserSponsor(t, "test", true)
			},
		},
		{
			"TestWithdrawSponsor",
			func(t *testing.T) {
				testWithDrawSponsor(t, "test", true, big.NewFloat(1))
			},
		},
		{
			"TestLockUserBalance1",
			func(t *testing.T) {
				testLockUserBalance(t, "test", mockUserOpHash, true, big.NewFloat(1))
			},
		},
		{
			"ReleaseUserOpHashLockWhenFail",
			func(t *testing.T) {
				testReleaseUserOpHashLockWhenFail(t, "test", mockUserOpHash, true)
			},
		},
		{
			"TestLockUserBalance2",
			func(t *testing.T) {
				testLockUserBalance(t, "test", mockUserOpHash2, true, big.NewFloat(1))
			},
		},
		{
			"ReleaseBalanceWithActualCost",
			func(t *testing.T) {
				testReleaseBalanceWithActualCost(t, "test", mockUserOpHash2, true, big.NewFloat(0.5))
			},
		},
		{
			"TestLockBalanceRelease",
			func(t *testing.T) {
				ReleaseExpireLockBalance()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
func testReleaseBalanceWithActualCost(t *testing.T, userId string, userOpHash []byte, isTestNet bool, actualCost *big.Float) {
	_, err := ReleaseBalanceWithActualCost(userId, userOpHash, actualCost, isTestNet)
	if err != nil {
		t.Error(err)
	}

}
func testReleaseUserOpHashLockWhenFail(t *testing.T, userId string, userOpHash []byte, isTestNet bool) {
	_, err := ReleaseUserOpHashLockWhenFail(userOpHash, isTestNet)
	if err != nil {
		t.Error(err)
	}
}
func testLockUserBalance(t *testing.T, userId string, userOpHash []byte, isTestNet bool, lockAmount *big.Float) {
	_, err := LockUserBalance(userId, userOpHash, isTestNet, lockAmount)
	if err != nil {
		t.Error(err)
	}
}
func testWithDrawSponsor(t *testing.T, userId string, isTestNet bool, amount *big.Float) {

	request := model.WithdrawSponsorRequest{
		PayUserId: userId,
		IsTestNet: isTestNet,
	}
	txHash := "0x123456"
	result, err := WithDrawSponsor(&request, txHash)
	if err != nil {
		t.Error(err)
	}
	if result == nil {
		t.Error("WithDrawSponsor failed")
	}
}
func testFindUserSponsor(t *testing.T, userId string, isTestNet bool) {
	sponsorModel, err := findUserSponsor(userId, isTestNet)
	if err != nil {
		t.Error(err)
	}
	if sponsorModel == nil {
		t.Error("findUserSponsor failed")
	}
	jsonSt, _ := json.Marshal(sponsorModel)
	t.Log(string(jsonSt))
}
