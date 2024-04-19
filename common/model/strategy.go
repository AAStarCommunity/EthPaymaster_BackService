package model

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
)

type Strategy struct {
	Id                 string                     `json:"id"`
	StrategyCode       string                     `json:"strategy_code"`
	PaymasterInfo      *PaymasterInfo             `json:"paymaster_info"`
	NetWorkInfo        *NetWorkInfo               `json:"network_info"`
	EntryPointInfo     *EntryPointInfo            `json:"entrypoint_info"`
	Description        string                     `json:"description"`
	ExecuteRestriction StrategyExecuteRestriction `json:"execute_restriction"`
}
type PaymasterInfo struct {
	PayMasterAddress *common.Address `json:"paymaster_address"`
	PayType          types.PayType   `json:"pay_type"`
}
type NetWorkInfo struct {
	NetWork types.Network   `json:"network"`
	Token   types.TokenType `json:"tokens"`
}
type EntryPointInfo struct {
	EntryPointAddress *common.Address         `json:"entrypoint_address"`
	EntryPointTag     types.EntrypointVersion `json:"entrypoint_tag"`
}

func (strategy *Strategy) GetPaymasterAddress() *common.Address {
	return strategy.PaymasterInfo.PayMasterAddress
}
func (strategy *Strategy) GetEntryPointAddress() *common.Address {
	return strategy.EntryPointInfo.EntryPointAddress
}
func (strategy *Strategy) GetNewWork() types.Network {
	return strategy.NetWorkInfo.NetWork
}

func (strategy *Strategy) GetUseToken() types.TokenType {
	return strategy.NetWorkInfo.Token
}
func (strategy *Strategy) GetPayType() types.PayType {
	return strategy.PaymasterInfo.PayType
}
func (strategy *Strategy) GetStrategyEntryPointTag() types.EntrypointVersion {
	return strategy.EntryPointInfo.EntryPointTag
}

type StrategyExecuteRestriction struct {
	BanSenderAddress   string             `json:"ban_sender_address"`
	EffectiveStartTime int64              `json:"effective_start_time"`
	EffectiveEndTime   int64              `json:"effective_end_time"`
	GlobalMaxUSD       int64              `json:"global_max_usd"`
	GlobalMaxOpCount   int64              `json:"global_max_op_count"`
	DayMaxUSD          int64              `json:"day_max_usd"`
	StartTime          int64              `json:"start_time"`
	EndTime            int64              `json:"end_time"`
	AccessProject      mapset.Set[string] `json:"access_project"`
}

type StrategyValidateConfig struct {
	ValidateContractAddress string `json:"validate_contract_address"`
}
