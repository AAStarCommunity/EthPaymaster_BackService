package model

import "math/big"

type DepositSponsorRequest struct {
	TimeStamp      int64  `json:"time_stamp"`
	DepositAddress string `json:"deposit_address"`
	TxHash         string `json:"tx_hash"`
	IsTestNet      bool   `json:"is_test_net"`
	PayUserId      string `json:"pay_user_id"`
	DepositSource  string `json:"deposit_source"`
}
type WithdrawSponsorRequest struct {
	Amount *big.Float

	PayUserId string
	IsTestNet bool
	TxInfo    map[string]string
	TxHash    string
}
type GetSponsorTransactionsRequest struct {
}
type GetSponsorMetaDataRequest struct {
}

type Transaction struct {
}
