package global_const

type PayType string

const (
	PayTypeProjectSponsor PayType = "PayTypeProjectSponsor"
	PayTypeERC20          PayType = "PayTypeERC20"
	PayTypeSuperVerifying PayType = "PayTypeSuperVerifying"
	PayTypeUserSponsor    PayType = "PayTypeUserSponsor"
)
