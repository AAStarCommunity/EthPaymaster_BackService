package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
)

func GetSupportStrategyExecute(network string) ([]model.Strategy, error) {
	return dashboard_service.GetStrategyListByNetwork(types.Network(network)), nil
}
