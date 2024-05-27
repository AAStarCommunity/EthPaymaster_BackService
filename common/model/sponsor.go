package model

import "math/big"

type DepositSponsorRequest struct {
	Source string
	Amount *big.Float
	TxHash string

	TxInfo    map[string]string
	PayUserId string
	IsTestNet bool
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
