package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
)

func GetSupportStrategyExecute(network string) ([]model.Strategy, error) {
	return dashboard_service.GetStrategyListByNetwork(global_const.Network(network)), nil
}
