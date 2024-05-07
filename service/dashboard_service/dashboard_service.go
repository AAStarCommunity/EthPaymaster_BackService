package dashboard_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"errors"
)

func GetStrategyByCode(strategyCode string) *model.Strategy {
	return config.GetBasicStrategyConfig(global_const.BasicStrategyCode(strategyCode))
}

func GetSuitableStrategy(entryPointVersion global_const.EntrypointVersion, chain global_const.Network, payType global_const.PayType) (*model.Strategy, error) {
	if entryPointVersion == "" {
		entryPointVersion = global_const.EntrypointV06
	}
	strategy, err := config.GetSuitableStrategy(entryPointVersion, chain, payType)
	if err != nil {
		return nil, err
	}

	if strategy == nil {
		return nil, errors.New("strategy not found")
	}
	return strategy, nil
}

func IsEntryPointsSupport(address string, chain global_const.Network) bool {
	supportEntryPointSet, _ := config.GetSupportEntryPoints(chain)
	if supportEntryPointSet == nil {
		return false
	}
	return supportEntryPointSet.Contains(address)
}
func IsPayMasterSupport(address string, chain global_const.Network) bool {
	supportPayMasterSet, _ := config.GetSupportPaymaster(chain)
	if supportPayMasterSet == nil {
		return false
	}

	return supportPayMasterSet.Contains(address)
}
