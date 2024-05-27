package sponsor_manager

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"encoding/json"
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
	changeModel := &UserSponsorBalanceUpdateLogDBModel{
		UserOpHash: utils.EncodeToHexStringWithPrefix(userOpHash),
		PayUserId:  userId,
		Amount:     lockAmount,
		Source:     "GasTank",
		IsTestNet:  isTestNet,
		UpdateType: global_const.UpdateTypeLock,
	}
	err = AddBalanceChangeLog(changeModel)
	if err != nil {
		return err
	}

	return nil
}

func ReleaseBalanceWithActualCost(userId string, userOpHash []byte, network global_const.Network,
	actualGasCost *big.Float, isTestNet bool) error {
	userOpHashHex := utils.EncodeToHexStringWithPrefix(userOpHash)

	changeModel, err := GetChangeModel(global_const.UpdateTypeLock, "", userOpHashHex, isTestNet)
	if err != nil {
		return err
	}
	balanceModel, err := getUserSponsorBalance(changeModel.PayUserId, changeModel.IsTestNet)

	lockBalance := changeModel.Amount
	balanceModel.LockBalance = new(big.Float).Sub(balanceModel.LockBalance, lockBalance)
	refundBalance := new(big.Float).Sub(lockBalance, actualGasCost)
	balanceModel.AvailableBalance = new(big.Float).Add(balanceModel.AvailableBalance, refundBalance)

	err = UpdateSponsor(balanceModel, isTestNet)

	if err != nil {
		return err
	}
	changeDBModel := &UserSponsorBalanceUpdateLogDBModel{
		UserOpHash: userOpHashHex,
		PayUserId:  changeModel.PayUserId,
		Amount:     refundBalance,
		Source:     "GasTank",
		IsTestNet:  isTestNet,
		UpdateType: global_const.UpdateTypeRelease,
	}
	err = AddBalanceChangeLog(changeDBModel)
	return nil
}

type ReleaseUserOpHashLockInput struct {
	UserOpHash []byte
}

func ReleaseUserOpHashLock(userOpHash []byte, isTestNet bool) (err error) {
	// Get ChangeLog By UserOpHash
	userOpHashHex := utils.EncodeToHexStringWithPrefix(userOpHash)
	changeModel, err := GetChangeModel(global_const.UpdateTypeLock, "", userOpHashHex, isTestNet)
	if err != nil {
		return err
	}
	// Release Lock
	balanceModel, err := getUserSponsorBalance(changeModel.PayUserId, changeModel.IsTestNet)

	lockBalance := changeModel.Amount

	balanceModel.LockBalance = new(big.Float).Sub(balanceModel.LockBalance, lockBalance)
	balanceModel.AvailableBalance = new(big.Float).Add(balanceModel.AvailableBalance, lockBalance)

	err = UpdateSponsor(balanceModel, isTestNet)
	if err != nil {
		return err
	}
	changeDBModel := &UserSponsorBalanceUpdateLogDBModel{
		UserOpHash: userOpHashHex,
		PayUserId:  changeModel.PayUserId,
		Amount:     lockBalance,
		Source:     "GasTank",
		IsTestNet:  isTestNet,
		UpdateType: global_const.UpdateTypeRelease,
	}
	err = AddBalanceChangeLog(changeDBModel)
	if err != nil {
		return err
	}
	return nil
}

//----------Functions----------

func DepositSponsor(input *model.DepositSponsorRequest) (*UserSponsorBalanceDBModel, error) {

	balanceModel, err := FindUserSponsorBalance(input.PayUserId, input.IsTestNet)

	tx := relayDB.Begin()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		balanceModel = &UserSponsorBalanceDBModel{}
		balanceModel.AvailableBalance = input.Amount
		balanceModel.PayUserId = input.PayUserId
		balanceModel.LockBalance = big.NewFloat(0)
		balanceModel.IsTestNet = input.IsTestNet
		err = tx.Create(balanceModel).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	newAvailableBalance := new(big.Float).Add(balanceModel.AvailableBalance, input.Amount)
	balanceModel.AvailableBalance = newAvailableBalance

	if updateErr := tx.Model(balanceModel).
		Where("pay_user_id = ?", balanceModel.PayUserId).
		Where("is_test_net = ?", input.IsTestNet).Updates(balanceModel).Error; updateErr != nil {
		tx.Rollback()
		return nil, updateErr
	}

	txInfoJSon, _ := json.Marshal(input.TxInfo)
	chagneModel := &UserSponsorBalanceUpdateLogDBModel{
		PayUserId:  input.PayUserId,
		Amount:     input.Amount,
		Source:     "Deposit",
		IsTestNet:  input.IsTestNet,
		UpdateType: global_const.UpdateTypeDeposit,
		TxHash:     input.TxHash,
		TxInfo:     txInfoJSon,
	}
	if createErr := tx.Create(chagneModel).Error; createErr != nil {
		tx.Rollback()
		return nil, createErr
	}

	tx.Commit()
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
	changeModel := &UserSponsorBalanceUpdateLogDBModel{
		PayUserId:  input.PayUserId,
		Amount:     input.Amount,
		Source:     "Withdraw",
		IsTestNet:  input.IsTestNet,
		UpdateType: global_const.UpdateTypeWithdraw,
		TxHash:     input.TxHash,
	}
	err = AddBalanceChangeLog(changeModel)
	if err != nil {
		return nil, err
	}
	return balanceModel, nil
}
