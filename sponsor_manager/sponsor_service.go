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

type UpdateType string
type BalanceType string

const (
	Deposit  UpdateType = "deposit"
	Lock     UpdateType = "lock"
	Withdraw UpdateType = "withdraw"
	Release  UpdateType = "release"

	AvailableBalance BalanceType = "available_balance"
	LockBalance      BalanceType = "lock_balance"
)

type UserSponsorBalanceDBModel struct {
	model.BaseData
	PayUserId        string     `gorm:"type:varchar(255);index" json:"pay_user_id"`
	AvailableBalance *big.Float `gorm:"type:numeric(30,18)" json:"available_balance"`
	LockBalance      *big.Float `gorm:"type:numeric(30,18)" json:"lock_balance"`
}

type DepositBalanceInput struct {
	Source    string
	Signature string
	Amount    *big.Float
	TxReceipt string
}

func (UserSponsorBalanceDBModel) TableName() string {
	return config.GetStrategyConfigTableName()
}

type UserSponsorBalanceUpdateLogDBModel struct {
	model.BaseData
	Amount     *big.Float `gorm:"type:numeric(30,18)" json:"amount"`
	UpdateType UpdateType `gorm:"type:varchar(20)" json:"update_type"`
}

func (UserSponsorBalanceUpdateLogDBModel) TableName() string {
	return config.GetStrategyConfigTableName()
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

func DepositSponsorBalance(input *DepositBalanceInput) (err error) {
	//TODO
	return nil

}

func LogBalanceChange(balanceType BalanceType, data interface{}, amount *big.Float) (err error) {
	//TODO
	return nil
}

func getUserSponsorBalance(userId string) (balanceModel *UserSponsorBalanceDBModel, err error) {
	relayDB.Model(&UserSponsorBalanceDBModel{}).Where("pay_user_id = ?", userId).First(&balanceModel)
	return balanceModel, nil
}
func CreateSponsorBalance(userId string) (err error) {
	//TODO
	return nil
}
