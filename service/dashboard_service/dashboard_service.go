package dashboard_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"errors"
)

// TODO just Temp Mock
var MockStrategyMap = map[string]*model.Strategy{}
var payMasterSupport = map[string]bool{}
var entryPointSupport = map[string]bool{}

func init() {

}
func GetStrategyById(strategyId string) *model.Strategy {
	return conf.BasicStrategyConfig[strategyId]
}

func GetSuitableStrategy(entrypoint string, chain types.Network, payType types.PayType) (*model.Strategy, error) {
	strategy := conf.SuitableStrategyMap[chain][entrypoint][payType]
	if strategy == nil {
		return nil, errors.New("strategy not found")
	}
	return strategy, nil
}
func IsEntryPointsSupport(address string) bool {
	if entryPointSupport[address] {
		return true
	}
	return false
}
func IsPayMasterSupport(address string) bool {
	if payMasterSupport[address] {
		return true
	}
	return false
}
