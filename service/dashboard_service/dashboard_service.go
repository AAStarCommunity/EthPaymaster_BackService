package dashboard_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"errors"
)

var MockStrategyMap = map[string]*model.Strategy{}

func init() {

}
func GetStrategyById(strategyId string) *model.Strategy {
	return conf.GetBasicStrategyConfig(strategyId)

}

func GetSuitableStrategy(entrypoint string, chain types.Network, payType types.PayType) (*model.Strategy, error) {
	strategy, err := conf.GetSuitableStrategy(entrypoint, chain, payType)
	if err != nil {
		return nil, err
	}

	if strategy == nil {
		return nil, errors.New("strategy not found")
	}
	return strategy, nil
}
func GetStrategyListByNetwork(chain types.Network) []model.Strategy {
	panic("implement me")
}
func IsEntryPointsSupport(address string, chain types.Network) bool {
	supportEntryPointSet, _ := conf.GetSupportEntryPoints(chain)
	if supportEntryPointSet == nil {
		return false
	}
	return supportEntryPointSet.Contains(address)
}
func IsPayMasterSupport(address string, chain types.Network) bool {
	supportPayMasterSet, _ := conf.GetSupportPaymaster(chain)
	if supportPayMasterSet == nil {
		return false
	}

	return supportPayMasterSet.Contains(address)
}
