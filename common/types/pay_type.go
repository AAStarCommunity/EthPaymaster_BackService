package types

type PayType string

const (
	PayTypeVerifying PayType = "0"
	PayTypeERC20     PayType = "1"
)
