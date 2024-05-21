package model

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Strategy struct {
	Id                 string                      `json:"id"`
	StrategyCode       string                      `json:"strategy_code"`
	PaymasterInfo      *PaymasterInfo              `json:"paymaster_info"`
	NetWorkInfo        *NetWorkInfo                `json:"network_info"`
	EntryPointInfo     *EntryPointInfo             `json:"entrypoint_info"`
	Description        string                      `json:"description"`
	ExecuteRestriction *StrategyExecuteRestriction `json:"execute_restriction"`
	Erc20TokenType     global_const.TokenType
}
type PaymasterInfo struct {
	PayMasterAddress        *common.Address      `json:"paymaster_address"`
	PayType                 global_const.PayType `json:"pay_type"`
	IsProjectErc20PayEnable bool                 `json:"is_project_erc20_pay_enable"`
}
type NetWorkInfo struct {
	NetWork  global_const.Network   `json:"network"`
	GasToken global_const.TokenType `json:"tokens"`
}
type EntryPointInfo struct {
	EntryPointAddress *common.Address                `json:"entrypoint_address"`
	EntryPointVersion global_const.EntrypointVersion `json:"entrypoint_version"`
}

func (strategy *Strategy) GetPaymasterAddress() *common.Address {
	return strategy.PaymasterInfo.PayMasterAddress
}
func (strategy *Strategy) GetEntryPointAddress() *common.Address {
	return strategy.EntryPointInfo.EntryPointAddress
}
func (strategy *Strategy) GetNewWork() global_const.Network {
	return strategy.NetWorkInfo.NetWork
}

func (strategy *Strategy) GetUseToken() global_const.TokenType {
	return strategy.NetWorkInfo.GasToken
}
func (strategy *Strategy) GetPayType() global_const.PayType {
	return strategy.PaymasterInfo.PayType
}
func (strategy *Strategy) GetStrategyEntrypointVersion() global_const.EntrypointVersion {
	return strategy.EntryPointInfo.EntryPointVersion
}
func (strategy *Strategy) IsCurrencyPayEnable() bool {
	return false
}

type StrategyExecuteRestriction struct {
	BanSenderAddress   mapset.Set[string] `json:"ban_sender_address"`
	EffectiveStartTime *big.Int           `json:"start_time"`
	EffectiveEndTime   *big.Int           `json:"end_time"`
	GlobalMaxUSD       *big.Float         `json:"global_max_usd"`
	GlobalMaxOpCount   *big.Int           `json:"global_max_op_count"`
	DayMaxUSD          *big.Float         `json:"day_max_usd"`
	AccessProject      mapset.Set[string] `json:"access_project"`
	AccessErc20        mapset.Set[string] `json:"access_erc20"`
	ChainIdWhiteList   mapset.Set[string] `json:"chain_id_whitelist"`
	Status             global_const.StrategyStatus
}

type StrategyValidateConfig struct {
	ValidateContractAddress string `json:"validate_contract_address"`
}
