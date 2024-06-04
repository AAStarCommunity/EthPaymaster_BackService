package sponsor_manager

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"gorm.io/datatypes"
	"math/big"
)

type UserSponsorBalanceUpdateLogDBModel struct {
	model.BaseData
	PayUserId  string                  `gorm:"type:varchar(255);index" json:"pay_user_id"`
	Amount     BigFloat                `gorm:"type:numeric(30,18)" json:"amount"`
	UpdateType global_const.UpdateType `gorm:"type:varchar(20)" json:"update_type"`
	UserOpHash string                  `gorm:"type:varchar(255)" json:"user_op_hash"`
	TxHash     string                  `gorm:"type:varchar(255)" json:"tx_hash"`
	TxInfo     datatypes.JSON          `gorm:"type:json" json:"tx_info"`
	Source     string                  `gorm:"type:varchar(255)" json:"source"`
	IsTestNet  bool                    `gorm:"type:boolean" json:"is_test_net"`
}

func (UserSponsorBalanceUpdateLogDBModel) TableName() string {
	return "relay_user_sponsor_balance_update_log"
}

func AddBalanceChangeLog(changeDbModel *UserSponsorBalanceUpdateLogDBModel) error {
	return relayDB.Create(changeDbModel).Error
}
func LogBalanceChange(updateType global_const.UpdateType, balanceType global_const.BalanceType, data interface{}, amount *big.Float) {

	//TODO
	return
}
func GetDepositAndWithDrawLog(userId string, IsTestNet bool) (models []*UserSponsorBalanceUpdateLogDBModel, err error) {
	tx := relayDB.Model(&UserSponsorBalanceUpdateLogDBModel{}).Where("pay_user_id = ?", userId).Where("is_test_net = ?", IsTestNet).Where("update_type in (?)", []global_const.UpdateType{global_const.UpdateTypeDeposit, global_const.UpdateTypeWithdraw}).Find(&models)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return models, nil
}
func LockBalanceChangeLog(payUserid string, userOpHash string, amount *big.Float, isTestNet bool, updateReason string) error {
	logModel := &UserSponsorBalanceUpdateLogDBModel{
		PayUserId:  payUserid,
		Amount:     BigFloat{amount},
		UserOpHash: userOpHash,
		Source:     "GasTank",
		IsTestNet:  isTestNet,
	}
	return relayDB.Create(logModel).Error
}

func GetChangeModel(updateType global_const.UpdateType, payUserId string, txHash string, isTestNet bool) (ChangeModel *UserSponsorBalanceUpdateLogDBModel, err error) {
	if updateType == global_const.UpdateTypeDeposit || updateType == global_const.UpdateTypeWithdraw {

		tx := relayDB.Model(ChangeModel).Where("tx_hash = ?", txHash).Where("pay_user_id", payUserId).Where("update_type = ?", global_const.UpdateTypeDeposit).Where("is_test_net", isTestNet).First(&ChangeModel)
		if tx.Error != nil {
			return nil, tx.Error
		} else {
			return ChangeModel, nil
		}
	} else if updateType == global_const.UpdateTypeLock || updateType == global_const.UpdateTypeRelease {
		tx := relayDB.Model(ChangeModel).Where("user_op_hash = ?", txHash).Where("update_type = ?", global_const.UpdateTypeLock).Where("is_test_net", isTestNet).First(&ChangeModel)
		if tx.Error != nil {
			return nil, tx.Error
		} else {
			return ChangeModel, nil
		}
	} else {
		return nil, nil
	}
}
