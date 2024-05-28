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

type Source string

const (
	SourceDashboard Source = "Dashboard"
	SourceRacks     Source = "Racks"
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
	balanceModel, err := findUserSponsor(userId, isTestNet)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerrors.Errorf("No Balance Deposit Here ")
		}
		return nil, err
	}
	return balanceModel.AvailableBalance.Float, nil
}

// LockUserBalance
// Reduce AvailableBalance and Increase LockBalance
func LockUserBalance(userId string, userOpHash []byte, isTestNet bool,
	lockAmount *big.Float) (*UserSponsorBalanceUpdateLogDBModel, error) {
	balanceModel, err := findUserSponsor(userId, isTestNet)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, xerrors.Errorf("No Balance Deposit Here ")
	}
	if err != nil {
		return nil, err
	}
	UserOphashStr := utils.EncodeToHexStringWithPrefix(userOpHash)
	_, err = GetChangeModel(global_const.UpdateTypeLock, userId, UserOphashStr, isTestNet)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, xerrors.Errorf("UserOpHash [%s] Already Lock", UserOphashStr)
	}

	if balanceModel.AvailableBalance.Cmp(lockAmount) < 0 {
		return nil, xerrors.Errorf("Insufficient balance [%s] not Enough to Lock [%s]", balanceModel.AvailableBalance.String(), lockAmount.String())
	}

	lockBalance := new(big.Float).Add(balanceModel.LockBalance.Float, lockAmount)
	availableBalance := new(big.Float).Sub(balanceModel.AvailableBalance.Float, lockAmount)
	balanceModel.LockBalance = BigFloat{lockBalance}
	balanceModel.AvailableBalance = BigFloat{availableBalance}
	tx := relayDB.Begin()
	if updateErr := tx.Model(&UserSponsorBalanceDBModel{}).
		Where("pay_user_id = ?", balanceModel.PayUserId).
		Where("is_test_net = ?", isTestNet).Updates(balanceModel).Error; updateErr != nil {
		tx.Rollback()
		return nil, updateErr
	}

	changeModel := &UserSponsorBalanceUpdateLogDBModel{
		UserOpHash: UserOphashStr,
		PayUserId:  userId,
		Amount:     BigFloat{lockAmount},
		IsTestNet:  isTestNet,
		UpdateType: global_const.UpdateTypeLock,
	}
	if createErr := tx.Create(changeModel).Error; createErr != nil {
		tx.Rollback()
		return nil, createErr
	}
	if commitErr := tx.Commit().Error; commitErr != nil {
		tx.Rollback()
		return nil, commitErr
	}

	return nil, nil
}

func ReleaseBalanceWithActualCost(userId string, userOpHash []byte,
	actualGasCost *big.Float, isTestNet bool) (*UserSponsorBalanceDBModel, error) {
	userOpHashHex := utils.EncodeToHexStringWithPrefix(userOpHash)

	changeModel, err := GetChangeModel(global_const.UpdateTypeLock, userId, userOpHashHex, isTestNet)
	if err != nil {
		return nil, err
	}
	balanceModel, err := findUserSponsor(changeModel.PayUserId, changeModel.IsTestNet)

	lockBalance := changeModel.Amount
	balanceModel.LockBalance = BigFloat{new(big.Float).Sub(balanceModel.LockBalance.Float, lockBalance.Float)}
	refundBalance := new(big.Float).Sub(lockBalance.Float, actualGasCost)
	balanceModel.AvailableBalance = BigFloat{new(big.Float).Add(balanceModel.AvailableBalance.Float, refundBalance)}

	tx := relayDB.Begin()
	if updateErr := tx.Model(&UserSponsorBalanceDBModel{}).
		Model(&UserSponsorBalanceDBModel{}).
		Where("pay_user_id = ?", balanceModel.PayUserId).
		Where("is_test_net = ?", isTestNet).Updates(balanceModel).Error; updateErr != nil {
		tx.Rollback()
		return nil, err
	}

	changeDBModel := &UserSponsorBalanceUpdateLogDBModel{
		UserOpHash: userOpHashHex,
		PayUserId:  changeModel.PayUserId,
		Amount:     BigFloat{refundBalance},
		Source:     "GasTank",
		IsTestNet:  isTestNet,
		UpdateType: global_const.UpdateTypeRelease,
	}
	if createErr := tx.Create(changeDBModel).Error; createErr != nil {
		tx.Rollback()
		return nil, createErr
	}
	if commitErr := tx.Commit().Error; commitErr != nil {
		tx.Rollback()
		return nil, commitErr
	}
	return balanceModel, nil
}

type ReleaseUserOpHashLockInput struct {
	UserOpHash []byte
}

func ReleaseUserOpHashLockWhenFail(userOpHash []byte, isTestNet bool) (*UserSponsorBalanceDBModel, error) {
	// Get ChangeLog By UserOpHash
	userOpHashHex := utils.EncodeToHexStringWithPrefix(userOpHash)
	changeModel, err := GetChangeModel(global_const.UpdateTypeLock, "", userOpHashHex, isTestNet)
	if err != nil {
		return nil, err
	}
	// Release Lock
	balanceModel, err := findUserSponsor(changeModel.PayUserId, changeModel.IsTestNet)
	lockBalance := changeModel.Amount
	balanceModel.LockBalance = BigFloat{new(big.Float).Sub(balanceModel.LockBalance.Float, lockBalance.Float)}
	balanceModel.AvailableBalance = BigFloat{new(big.Float).Add(balanceModel.AvailableBalance.Float, lockBalance.Float)}
	tx := relayDB.Begin()
	if updateErr := tx.Model(&UserSponsorBalanceDBModel{}).
		Where("pay_user_id = ?", balanceModel.PayUserId).
		Where("is_test_net = ?", isTestNet).Updates(balanceModel).Error; updateErr != nil {
		tx.Rollback()
		return nil, err
	}

	changeDBModel := &UserSponsorBalanceUpdateLogDBModel{
		UserOpHash: userOpHashHex,
		PayUserId:  changeModel.PayUserId,
		Amount:     lockBalance,
		IsTestNet:  isTestNet,
		UpdateType: global_const.UpdateTypeRelease,
	}
	if createErr := tx.Create(changeDBModel).Error; createErr != nil {
		tx.Rollback()
		return nil, createErr
	}
	if commitErr := tx.Commit().Error; commitErr != nil {
		tx.Rollback()
		return nil, commitErr
	}
	return nil, err
}

//----------Functions----------

func DepositSponsor(input *model.DepositSponsorRequest) (*UserSponsorBalanceDBModel, error) {

	balanceModel, err := FindUserSponsorBalance(input.PayUserId, input.IsTestNet)

	tx := relayDB.Begin()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		//init Data
		balanceModel = &UserSponsorBalanceDBModel{}
		balanceModel.AvailableBalance = BigFloat{big.NewFloat(0)}
		balanceModel.PayUserId = input.PayUserId
		balanceModel.LockBalance = BigFloat{big.NewFloat(0)}
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
	newAvailableBalance := BigFloat{new(big.Float).Add(balanceModel.AvailableBalance.Float, input.Amount)}
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
		Amount:     BigFloat{input.Amount},
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

	if commitErr := tx.Commit().Error; commitErr != nil {
		tx.Rollback()
		return nil, commitErr
	}
	return balanceModel, nil
}

func WithDrawSponsor(input *model.WithdrawSponsorRequest) (*UserSponsorBalanceDBModel, error) {
	balanceModel, err := FindUserSponsorBalance(input.PayUserId, input.IsTestNet)
	if err != nil {
		return nil, err
	}
	if balanceModel.AvailableBalance.Cmp(input.Amount) < 0 {
		return nil, xerrors.Errorf("Insufficient balance [%s] not Enough to Withdraw [%s]", balanceModel.AvailableBalance.String(), input.Amount.String())
	}
	newAvailableBalance := new(big.Float).Sub(balanceModel.AvailableBalance.Float, input.Amount)
	balanceModel.AvailableBalance = BigFloat{newAvailableBalance}
	tx := relayDB.Begin()
	if updateErr := tx.Model(&UserSponsorBalanceDBModel{}).
		Where("pay_user_id = ?", balanceModel.PayUserId).
		Where("is_test_net = ?", input.IsTestNet).Updates(balanceModel).Error; updateErr != nil {
		tx.Rollback()
		return nil, updateErr
	}
	changeModel := &UserSponsorBalanceUpdateLogDBModel{
		PayUserId:  input.PayUserId,
		Amount:     BigFloat{input.Amount},
		Source:     "Withdraw",
		IsTestNet:  input.IsTestNet,
		UpdateType: global_const.UpdateTypeWithdraw,
		TxHash:     input.TxHash,
	}
	if createErr := tx.Create(changeModel).Error; createErr != nil {
		tx.Rollback()
		return nil, createErr
	}
	if commitErr := tx.Commit().Error; commitErr != nil {
		tx.Rollback()
		return nil, commitErr
	}
	return balanceModel, nil
}
