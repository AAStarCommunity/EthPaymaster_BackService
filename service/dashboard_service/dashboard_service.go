package dashboard_service

import "AAStarCommunity/EthPaymaster_BackService/common/model"

var mockStrategyMap = map[string]model.Strategy{}

func init() {
	mockStrategyMap["1"] = model.Strategy{
		Id:                "1",
		EntryPointAddress: "0x123",
		PayMasterAddress:  "0x123",
	}
}
func GetStrategyById(strategyId string) model.Strategy {
	return mockStrategyMap[strategyId]
}
