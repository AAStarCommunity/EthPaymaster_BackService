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

func GetAvailableBalance(userId string) (balance *big.Float, err error) {
	balanceModel, err := getUserSponsorBalance(userId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerrors.Errorf("No Balance Deposit Here ")
		}
		return nil, err
	}
	return balanceModel.AvailableBalance, nil
}

func LockUserBalance(userId string, userOpHash []byte, network global_const.Network,
	lockAmount *big.Float) (err error) {
	balanceModel, err := getUserSponsorBalance(userId)
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
	tx := relayDB.Model(&UserSponsorBalanceDBModel{}).Where("pay_user_id = ?", userId).Updates(balanceModel)
	if tx.Error != nil {
		return tx.Error
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
	//TODO
	return nil
}

func getUserSponsorBalance(userId string) (balanceModel *UserSponsorBalanceDBModel, err error) {
	tx := relayDB.Model(&UserSponsorBalanceDBModel{}).Where("pay_user_id = ?", userId).First(&balanceModel)
	return balanceModel, tx.Error
}

//----------Functions----------

func DepositSponsor(input *model.DepositSponsorRequest) (balanceModel *UserSponsorBalanceDBModel, err error) {
	//TODO
	balanceModel, err = FindUserSponsorBalance(input.PayUserId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		balanceModel.AvailableBalance = input.Amount
		balanceModel.PayUserId = input.PayUserId
		balanceModel.LockBalance = big.NewFloat(0)
		tx := relayDB.Create(&balanceModel)
		if tx.Error != nil {
			return nil, tx.Error
		}
	}
	if err != nil {
		return nil, err
	}
	newAvailableBalance := new(big.Float).Add(balanceModel.AvailableBalance, input.Amount)
	balanceModel.AvailableBalance = newAvailableBalance
	tx := relayDB.Model(&UserSponsorBalanceDBModel{}).Where("pay_user_id = ?", input.PayUserId).Updates(balanceModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	LogBalanceChange(global_const.UpdateTypeDeposit, global_const.BalanceTypeAvailableBalance, input, input.Amount)
	return balanceModel, nil
}

func WithDrawSponsor(input *model.WithdrawSponsorRequest) (balanceModel *UserSponsorBalanceDBModel, err error) {
	balanceModel, err = FindUserSponsorBalance(input.PayUserId)
	if err != nil {
		return nil, err
	}
	if balanceModel.AvailableBalance.Cmp(input.Amount) < 0 {
		return nil, xerrors.Errorf("Insufficient balance [%s] not Enough to Withdraw [%s]", balanceModel.AvailableBalance.String(), input.Amount.String())
	}
	newAvailableBalance := new(big.Float).Sub(balanceModel.AvailableBalance, input.Amount)
	balanceModel.AvailableBalance = newAvailableBalance
	err = UpdateSponsorBalance(balanceModel)
	if err != nil {
		return nil, err
	}
	LogBalanceChange(global_const.UpdateTypeWithdraw, global_const.BalanceTypeAvailableBalance, input, input.Amount)
	return balanceModel, nil
}
func GetSponsorTransactionList() (err error) {
	//TODO
	return nil
}

func GetSponsorMetaData() (err error) {
	//TODO
	return nil
}
