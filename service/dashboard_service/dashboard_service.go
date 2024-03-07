package dashboard_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"errors"
)

var mockStrategyMap = map[string]*model.Strategy{}

func init() {
	mockStrategyMap["1"] = &model.Strategy{
		Id:                "1",
		EntryPointAddress: "0x0576a174D229E3cFA37253523E645A78A0C91B57",
		PayMasterAddress:  "0x0000000000325602a77416A16136FDafd04b299f",
		NetWork:           types.Sepolia,
		Token:             types.USDT,
	}
}
func GetStrategyById(strategyId string) *model.Strategy {
	return mockStrategyMap[strategyId]
}

func GetSuitableStrategy(entrypoint string, chain types.NetWork, token string) (*model.Strategy, error) {
	return nil, errors.New("not implemented")
}
