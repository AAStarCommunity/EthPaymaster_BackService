package model

type DepositSponsorRequest struct {
	TimeStamp     int64  `json:"time_stamp"`
	TxHash        string `json:"tx_hash"`
	IsTestNet     bool   `json:"is_test_net"`
	PayUserId     string `json:"pay_user_id"`
	DepositSource string `json:"deposit_source"`
}
type WithdrawSponsorRequest struct {
	Amount         float64 `json:"amount"`
	TimeStamp      int64   `json:"time_stamp"`
	PayUserId      string  `json:"pay_user_id"`
	IsTestNet      bool    `json:"is_test_net"`
	WithdrawSource string  `json:"withdraw_source"`
	RefundAddress  string  `json:"refund_address"`
	DepositSource  string  `json:"deposit_source"`
}
