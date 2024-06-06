package model

import "math/big"

type DepositSponsorRequest struct {
	Source string     `json:"source"`
	Amount *big.Float `json:"amount"`
	TxHash string     `json:"tx_hash"`

	TxInfo    map[string]string `json:"tx_info"`
	PayUserId string            `json:"pay_user_id"`
	IsTestNet bool              `json:"is_test_net"`
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
