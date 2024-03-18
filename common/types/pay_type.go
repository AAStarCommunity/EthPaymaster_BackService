package types

type PayType string

const (
	PayTypeVerifying PayType = "00"
	PayTypeERC20     PayType = "01"
)
