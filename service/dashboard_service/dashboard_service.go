package dashboard_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"errors"
)

func GetStrategyByCode(strategyCode string) *model.Strategy {
	return conf.GetBasicStrategyConfig(global_const.BasicStrategyCode(strategyCode))

}

func GetSuitableStrategy(entrypoint string, chain global_const.Network, payType global_const.PayType) (*model.Strategy, error) {
	strategy, err := conf.GetSuitableStrategy(entrypoint, chain, payType)
	if err != nil {
		return nil, err
	}

	if strategy == nil {
		return nil, errors.New("strategy not found")
	}
	return strategy, nil
}
func GetStrategyListByNetwork(chain global_const.Network) []model.Strategy {
	return nil
	//TODO
}
func IsEntryPointsSupport(address string, chain global_const.Network) bool {
	supportEntryPointSet, _ := conf.GetSupportEntryPoints(chain)
	if supportEntryPointSet == nil {
		return false
	}
	return supportEntryPointSet.Contains(address)
}
func IsPayMasterSupport(address string, chain global_const.Network) bool {
	supportPayMasterSet, _ := conf.GetSupportPaymaster(chain)
	if supportPayMasterSet == nil {
		return false
	}

	return supportPayMasterSet.Contains(address)
}
