package global_const

type PayType string

const (
	PayTypeVerifying      PayType = "00"
	PayTypeERC20          PayType = "01"
	PayTypeSuperVerifying PayType = "02"
)
