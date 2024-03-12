package model

import "AAStarCommunity/EthPaymaster_BackService/common/types"

type Strategy struct {
	Id                 string                     `json:"id"`
	EntryPointAddress  string                     `json:"entrypoint_address"`
	EntryPointTag      types.EntrypointTag        `json:"entrypoint_tag"`
	PayMasterAddress   string                     `json:"paymaster_address"`
	PayType            types.PayType              `json:"pay_type"`
	NetWork            types.Network              `json:"network"`
	Token              types.TokenType            `json:"token"`
	Description        string                     `json:"description"`
	ExecuteRestriction StrategyExecuteRestriction `json:"execute_restriction"`
	EnableEoa          bool                       `json:"enable_eoa"`
	Enable7560         bool                       `json:"enable_7560"`
	EnableErc20        bool                       `json:"enable_erc20"`
	Enable4844         bool                       `json:"enable_4844"`
	EnableCurrency     bool                       `json:"enable_currency"`
}
type StrategyExecuteRestriction struct {
	BanSenderAddress   string `json:"ban_sender_address"`
	EffectiveStartTime int64  `json:"effective_start_time"`
	EffectiveEndTime   int64  `json:"effective_end_time"`
	GlobalMaxUSD       int64  `json:"global_max_usd"`
	GlobalMaxOpCount   int64  `json:"global_max_op_count"`
	DayMaxUSD          int64  `json:"day_max_usd"`
}

type StrategyValidateConfig struct {
	ValidateContractAddress string `json:"validate_contract_address"`
}
