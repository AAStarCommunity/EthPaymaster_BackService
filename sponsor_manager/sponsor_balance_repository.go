package sponsor_manager

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/config"
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
	return config.GetStrategyConfigTableName()
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
