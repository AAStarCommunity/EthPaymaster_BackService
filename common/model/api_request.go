package model

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"errors"
	"golang.org/x/xerrors"
)

type TryPayUserOpRequest struct {
	ForceStrategyId        string         `json:"force_strategy_id"`
	ForceNetwork           types.Network  `json:"force_network"`
	ForceToken             string         `json:"force_token"`
	ForceEntryPointAddress string         `json:"force_entrypoint_address"`
	UserOp                 map[string]any `json:"user_operation"`
	Extra                  interface{}    `json:"extra"`
	OnlyEstimateGas        bool           `json:"only_estimate_gas"`
}

func (request *TryPayUserOpRequest) Validate() error {
	if len(request.ForceStrategyId) == 0 {
		if len(request.ForceNetwork) == 0 || len(request.ForceToken) == 0 || len(request.ForceEntryPointAddress) == 0 {
			return errors.New("strategy configuration illegal")
		}
	}
	if request.ForceStrategyId == "" && (request.ForceToken == "" || request.ForceNetwork == "") {
		return xerrors.Errorf("Token And Network Must Set When ForceStrategyId Is Empty")
	}
	if conf.Environment.IsDevelopment() && request.ForceNetwork != "" {
		//if types.TestNetWork[request.ForceNetwork] {
		//	return xerrors.Errorf(" %s not the Test Network ", request.ForceNetwork)
		//}
	}
	exist := conf.CheckEntryPointExist(request.ForceNetwork, request.ForceEntryPointAddress)
	if !exist {
		return xerrors.Errorf("ForceEntryPointAddress: [%s] not exist in [%s] network", request.ForceEntryPointAddress, request.ForceNetwork)
	}
	return nil
}
