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
}

func (UserSponsorBalanceDBModel) TableName() string {
	return config.GetStrategyConfigTableName()
}

func CreateSponsorBalance(balanceModel *UserSponsorBalanceDBModel) error {
	return relayDB.Create(balanceModel).Error
}
func FindUserSponsorBalance(userId string) (balanceModel *UserSponsorBalanceDBModel, err error) {
	balanceModel = &UserSponsorBalanceDBModel{}
	tx := relayDB.Where("pay_user_id = ?", userId).First(balanceModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return balanceModel, nil
}
func UpdateSponsorBalance(balanceModel *UserSponsorBalanceDBModel) error {
	tx := relayDB.Model(&UserSponsorBalanceDBModel{}).Where("pay_user_id = ?", balanceModel.PayUserId).Updates(balanceModel)
	return tx.Error
}
