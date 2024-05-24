package model

import "math/big"

type DepositSponsorRequest struct {
	Source    string
	Amount    *big.Float
	TxReceipt string
	PayUserId string
}
type WithdrawSponsorRequest struct {
	Source    string
	Amount    *big.Float
	TxReceipt string
	PayUserId string
}
type GetSponsorTransactionsRequest struct {
}
type GetSponsorMetaDataRequest struct {
}

type Transaction struct {
}
