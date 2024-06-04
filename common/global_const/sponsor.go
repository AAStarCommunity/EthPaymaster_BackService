package global_const

type UpdateType string

const (
	UpdateTypeDeposit  UpdateType = "deposit"
	UpdateTypeLock     UpdateType = "lock"
	UpdateTypeWithdraw UpdateType = "withdraw"
	UpdateTypeRelease  UpdateType = "release"
)

type BalanceType string

const (
	BalanceTypeAvailableBalance BalanceType = "available_balance"
	BalanceTypeLockBalance      BalanceType = "lock_balance"
)
