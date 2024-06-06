package sponsor_manager

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"database/sql/driver"
	"fmt"
	"math/big"
)

type UserSponsorBalanceDBModel struct {
	model.BaseData
	PayUserId        string   `gorm:"type:varchar(255);index" json:"pay_user_id"`
	AvailableBalance BigFloat `gorm:"type:numeric(30,18)" json:"available_balance"`
	SponsoredBalance BigFloat `gorm:"type:numeric(30,18)" json:"sponsored_balance"`
	LockBalance      BigFloat `gorm:"type:numeric(30,18)" json:"lock_balance"`
	Source           string   `gorm:"type:varchar(255)" json:"source"`
	SponsorAddress   string   `gorm:"type:varchar(255)" json:"sponsor_address"`
	IsTestNet        bool     `gorm:"type:boolean" json:"is_test_net"`
}

func (UserSponsorBalanceDBModel) TableName() string {
	return "relay_user_sponsor_balance"
}

func CreateSponsorBalance(balanceModel *UserSponsorBalanceDBModel) error {
	return relayDB.Create(balanceModel).Error
}
func FindUserSponsorBalance(userId string, isTestNet bool) (balanceModel *UserSponsorBalanceDBModel, err error) {

	balanceModel = &UserSponsorBalanceDBModel{}
	tx := relayDB.Where("pay_user_id = ?", userId).Where("is_test_net = ?", isTestNet).First(balanceModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return balanceModel, nil
}
func UpdateSponsor(balanceModel *UserSponsorBalanceDBModel, isTestNet bool) error {
	tx := relayDB.Model(&UserSponsorBalanceDBModel{}).Where("pay_user_id = ?", balanceModel.PayUserId).Where("is_test_net = ?", isTestNet).Updates(balanceModel)
	return tx.Error
}
func findUserSponsor(userId string, isTestNet bool) (balanceModel *UserSponsorBalanceDBModel, err error) {
	tx := relayDB.Model(&UserSponsorBalanceDBModel{}).Where("pay_user_id = ?", userId).Where("is_test_net = ?", isTestNet).First(&balanceModel)
	return balanceModel, tx.Error
}

// BigFloat wraps big.Float to implement sql.Scanner
type BigFloat struct {
	*big.Float
}

// Scan implements the sql.Scanner interface for BigFloat
func (bf *BigFloat) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		f, _, err := big.ParseFloat(string(v), 10, 256, big.ToNearestEven)
		if err != nil {
			return err
		}
		bf.Float = f
	case string:
		f, _, err := big.ParseFloat(v, 10, 256, big.ToNearestEven)
		if err != nil {
			return err
		}
		bf.Float = f
	case float64:
		bf.Float = big.NewFloat(v)
	case nil:
		bf.Float = nil
	default:
		return fmt.Errorf("cannot scan type %T into BigFloat", value)
	}
	return nil
}
func (bf *BigFloat) Value() (driver.Value, error) {
	if bf.Float == nil {
		return nil, nil
	}
	return bf.Text('f', -1), nil
}
