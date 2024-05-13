package model

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
)

type UserOpRequest struct {
	ForceStrategyId        string                         `json:"force_strategy_id"`
	ForceNetwork           global_const.Network           `json:"force_network"`
	Erc20Token             global_const.TokenType         `json:"force_token"`
	ForceEntryPointAddress string                         `json:"force_entrypoint_address"`
	UserOp                 map[string]any                 `json:"user_operation"`
	Extra                  interface{}                    `json:"extra"`
	EstimateOpGas          bool                           `json:"estimate_op_gas"`
	EntryPointVersion      global_const.EntrypointVersion `json:"entrypoint_version"`
}
type JsonRpcRequest struct {
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
}
type ClientCredential struct {
	ApiKey string `json:"apiKey"`
}
