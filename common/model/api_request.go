package model

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
)

type UserOpRequest struct {
	ForceStrategyId        string               `json:"force_strategy_id"`
	ForceNetwork           global_const.Network `json:"force_network"`
	Erc20Token             string               `json:"force_token"`
	ForceEntryPointAddress string               `json:"force_entrypoint_address"`
	UserOp                 map[string]any       `json:"user_operation"`
	Extra                  interface{}          `json:"extra"`
	EstimateOpGas          bool                 `json:"estimate_op_gas"`
}
