package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
)

func GetSupportStrategyExecute(network string) (map[string]*model.Strategy, error) {
	return dashboard_service.MockStrategyMap, nil
}
