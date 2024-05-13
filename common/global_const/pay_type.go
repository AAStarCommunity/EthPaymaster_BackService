package global_const

type PayType string

const (
	PayTypeVerifying      PayType = "PayTypeVerifying"
	PayTypeERC20          PayType = "PayTypeERC20"
	PayTypeSuperVerifying PayType = "PayTypeSuperVerifying"
)
