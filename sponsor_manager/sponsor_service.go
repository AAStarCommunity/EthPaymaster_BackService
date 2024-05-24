package sponsor_manager

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"errors"
	"golang.org/x/xerrors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math/big"
	"sync"
)

var (
	relayDB  *gorm.DB
	onlyOnce = sync.Once{}
)

func Init() {
	onlyOnce.Do(func() {
		relayDBDsn := config.GetRelayDBDSN()

		relayDBVar, err := gorm.Open(postgres.Open(relayDBDsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		relayDB = relayDBVar
	})
}

//----------Functions----------

func GetAvailableBalance(userId string, isTestNet bool) (balance *big.Float, err error) {
	balanceModel, err := getUserSponsorBalance(userId, isTestNet)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerrors.Errorf("No Balance Deposit Here ")
		}
		return nil, err
	}
	return balanceModel.AvailableBalance, nil
}

// LockUserBalance
// Reduce AvailableBalance and Increase LockBalance
func LockUserBalance(userId string, userOpHash []byte, isTestNet bool,
	lockAmount *big.Float) (err error) {
	balanceModel, err := getUserSponsorBalance(userId, isTestNet)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return xerrors.Errorf("No Balance Deposit Here ")
	}
	if err != nil {
		return err
	}
	lockBalance := new(big.Float).Add(balanceModel.LockBalance, lockAmount)
	availableBalance := new(big.Float).Sub(balanceModel.AvailableBalance, lockAmount)
	balanceModel.LockBalance = lockBalance
	balanceModel.AvailableBalance = availableBalance
	err = UpdateSponsor(balanceModel, isTestNet)
	if err != nil {
		return err
	}
	LogBalanceChange(global_const.UpdateTypeLock, global_const.BalanceTypeLockBalance, userOpHash, lockAmount)
	return nil
}

func ReleaseBalanceWithActualCost(userId string, userOpHash []byte, network global_const.Network,
	actualGasCost *big.Float) (err error) {
	//TODO
	return nil
}

type ReleaseUserOpHashLockInput struct {
	UserOpHash []byte
}

func ReleaseUserOpHashLock(userOpHash []byte) (err error) {
	// Get ChangeLog By UserOpHash
	return nil
}

//----------Functions----------

func DepositSponsor(input *model.DepositSponsorRequest) (balanceModel *UserSponsorBalanceDBModel, err error) {
	balanceModel, err = FindUserSponsorBalance(input.PayUserId, input.IsTestNet)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		balanceModel.AvailableBalance = input.Amount
		balanceModel.PayUserId = input.PayUserId
		balanceModel.LockBalance = big.NewFloat(0)
		balanceModel.IsTestNet = input.IsTestNet
		err = CreateSponsorBalance(balanceModel)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	newAvailableBalance := new(big.Float).Add(balanceModel.AvailableBalance, input.Amount)
	balanceModel.AvailableBalance = newAvailableBalance
	err = UpdateSponsor(balanceModel, input.IsTestNet)
	if err != nil {
		return nil, err
	}
	LogBalanceChange(global_const.UpdateTypeDeposit, global_const.BalanceTypeAvailableBalance, input, input.Amount)
	return balanceModel, nil
}

func WithDrawSponsor(input *model.WithdrawSponsorRequest) (balanceModel *UserSponsorBalanceDBModel, err error) {
	balanceModel, err = FindUserSponsorBalance(input.PayUserId, input.IsTestNet)
	if err != nil {
		return nil, err
	}
	if balanceModel.AvailableBalance.Cmp(input.Amount) < 0 {
		return nil, xerrors.Errorf("Insufficient balance [%s] not Enough to Withdraw [%s]", balanceModel.AvailableBalance.String(), input.Amount.String())
	}
	newAvailableBalance := new(big.Float).Sub(balanceModel.AvailableBalance, input.Amount)
	balanceModel.AvailableBalance = newAvailableBalance
	err = UpdateSponsor(balanceModel, input.IsTestNet)
	if err != nil {
		return nil, err
	}
	LogBalanceChange(global_const.UpdateTypeWithdraw, global_const.BalanceTypeAvailableBalance, input, input.Amount)
	return balanceModel, nil
}
