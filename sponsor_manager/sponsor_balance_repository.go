package sponsor_manager

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"fmt"
	"math/big"
)

type UserSponsorBalanceDBModel struct {
	model.BaseData
	PayUserId        string     `gorm:"type:varchar(255);index" json:"pay_user_id"`
	AvailableBalance *big.Float `gorm:"type:numeric(30,18)" json:"available_balance"`
	LockBalance      *big.Float `gorm:"type:numeric(30,18)" json:"lock_balance"`
	Source           string     `gorm:"type:varchar(255)" json:"source"`
	SponsorAddress   string     `gorm:"type:varchar(255)" json:"sponsor_address"`
	IsTestNet        bool       `gorm:"type:boolean" json:"is_test_net"`
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
func getUserSponsorBalance(userId string, isTestNet bool) (balanceModel *UserSponsorBalanceDBModel, err error) {
	tx := relayDB.Model(&UserSponsorBalanceDBModel{}).Where("pay_user_id = ?", userId).Where("is_test_net = ?", isTestNet).First(&balanceModel)
	return balanceModel, tx.Error
}

// BigFloat wraps big.Float to implement sql.Scanner
type BigFloat struct {
	*big.Float
}

// Scan implements the sql.Scanner interface for BigFloat
func (b *BigFloat) Scan(value interface{}) error {
	if b == nil {
		fmt.Println("BigFloat Scan value is nil")
		return nil
	}
	fmt.Println("BigFloat Scan value", value)
	switch v := value.(type) {
	case string:

		if _, ok := b.SetString(v); !ok {
			return fmt.Errorf("failed to parse string as *big.Float: %s", v)
		}
	case []byte:
		fmt.Println("BigFloat Scan value Case Byte", value)

		if _, ok := b.SetString(string(v)); !ok {
			return fmt.Errorf("failed to parse byte slice as *big.Float: %s", string(v))
		}
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}

func SelectUserSponsorBalanceDBModelWithScanList() (userSponsorBalanceDBModelWithScanList []UserSponsorBalanceDBModelWithScan, err error) {
	userSponsorBalanceDBModelWithScanList = make([]UserSponsorBalanceDBModelWithScan, 0)
	tx := relayDB.Model(&UserSponsorBalanceDBModelWithScan{}).Find(&userSponsorBalanceDBModelWithScanList)
	if tx.Error != nil {
		err = tx.Error
		return
	}
	fmt.Println(userSponsorBalanceDBModelWithScanList)
	return userSponsorBalanceDBModelWithScanList, nil
}

type UserSponsorBalanceDBModelWithScan struct {
	model.BaseData
	PayUserId        string  `gorm:"type:varchar(255);index" json:"pay_user_id"`
	AvailableBalance float64 `gorm:"type:numeric(30,18)" json:"available_balance"`
	LockBalance      float64 `gorm:"type:numeric(30,18)" json:"lock_balance"`
	Source           string  `gorm:"type:varchar(255)" json:"source"`
	SponsorAddress   string  `gorm:"type:varchar(255)" json:"sponsor_address"`
	IsTestNet        bool    `gorm:"type:boolean" json:"is_test_net"`
}

func (UserSponsorBalanceDBModelWithScan) TableName() string {
	return "relay_user_sponsor_balance"
}
