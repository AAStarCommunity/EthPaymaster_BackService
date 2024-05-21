package sponsor_manager

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"math/big"
)

func GetAvailableBalance(userId string) (balance *big.Float, err error) {
	// TODO
	return big.NewFloat(12.1), nil
}

func LockBalance(userId string, userOpHash []byte, network global_const.Network,
	lockAmount *big.Float) (err error) {
	//TODO
	return nil
}

func ReleaseBalanceWithActualCost(userId string, userOpHash []byte, network global_const.Network,
	actualGasCost *big.Float) (err error) {
	//TODO
	return nil
}
func ReleaseUserOpHashLock(userOpHash []byte) (err error) {
	//TODO
	return nil
}
