package sponsor_manager

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"gorm.io/datatypes"
	"math/big"
)

type UserSponsorBalanceUpdateLogDBModel struct {
	model.BaseData
	PayUserId  string                  `gorm:"type:varchar(255);index" json:"pay_user_id"`
	Amount     *big.Float              `gorm:"type:numeric(30,18)" json:"amount"`
	UpdateType global_const.UpdateType `gorm:"type:varchar(20)" json:"update_type"`
	UserOpHash string                  `gorm:"type:varchar(255)" json:"user_op_hash"`
	TxHash     string                  `gorm:"type:varchar(255)" json:"tx_hash"`
	TxInfo     datatypes.JSON          `gorm:"type:json" json:"tx_info"`
	Source     string                  `gorm:"type:varchar(255)" json:"source"`
	IsTestNet  bool                    `gorm:"type:boolean" json:"is_test_net"`
}

func (UserSponsorBalanceUpdateLogDBModel) TableName() string {
	return config.GetStrategyConfigTableName()
}

func LogBalanceChange(updateType global_const.UpdateType, balanceType global_const.BalanceType, data interface{}, amount *big.Float) {
	//TODO
	return
}
func GetDepositAndWithDrawLog(userId string, IsTestNet bool) (models []*UserSponsorBalanceUpdateLogDBModel, err error) {
	tx := relayDB.Model(&UserSponsorBalanceUpdateLogDBModel{}).Where("pay_user_id = ?", userId).Where("is_test_net = ?", IsTestNet).Where("update_type = ?", global_const.UpdateTypeDeposit).Or("update_type = ?", global_const.UpdateTypeWithdraw).Find(&models)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return models, nil
}
