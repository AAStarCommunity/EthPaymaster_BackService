package model

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
)

type UserOpRequest struct {
	StrategyCode      string                         `json:"strategy_code"`
	Network           global_const.Network           `json:"network"`
	UserPayErc20Token global_const.TokenType         `json:"user_pay_erc20_token"`
	UserOp            map[string]any                 `json:"user_operation"`
	Extra             interface{}                    `json:"extra"`
	EstimateOpGas     bool                           `json:"estimate_op_gas"`
	EntryPointVersion global_const.EntrypointVersion `json:"entrypoint_version"`
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
