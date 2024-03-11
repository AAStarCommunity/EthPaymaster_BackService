package dashboard_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"errors"
)

// TODO just Temp Mock
var mockStrategyMap = map[string]*model.Strategy{}
var payMasterSupport = map[string]bool{}
var entryPointSupport = map[string]bool{}

func init() {
	mockStrategyMap["1"] = &model.Strategy{
		Id:                "1",
		EntryPointAddress: "0x0576a174D229E3cFA37253523E645A78A0C91B57",
		PayMasterAddress:  "0x0000000000325602a77416A16136FDafd04b299f",
		NetWork:           types.Sepolia,
		Token:             types.USDT,
	}
	mockStrategyMap["2"] = &model.Strategy{
		Id:                "2",
		EntryPointAddress: "0x0576a174D229E3cFA37253523E645A78A0C91B57",
		PayMasterAddress:  "0x0000000000325602a77416A16136FDafd04b299f",
		NetWork:           types.Sepolia,
		Token:             types.ETH,
	}

	entryPointSupport["0x0576a174D229E3cFA37253523E645A78A0C91B57"] = true
	payMasterSupport["0x0000000000325602a77416A16136FDafd04b299f"] = true
}
func GetStrategyById(strategyId string) *model.Strategy {
	return mockStrategyMap[strategyId]
}
func GetSupportEntryPoint() {

}
func GetSuitableStrategy(entrypoint string, chain types.Network, token string) (*model.Strategy, error) {
	return nil, errors.New("not implemented")
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
