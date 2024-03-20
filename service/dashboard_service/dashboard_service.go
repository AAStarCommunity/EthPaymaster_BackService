package dashboard_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"errors"
)

// TODO just Temp Mock
var MockStrategyMap = map[string]*model.Strategy{}
var payMasterSupport = map[string]bool{}
var entryPointSupport = map[string]bool{}

func init() {
	MockStrategyMap["1"] = &model.Strategy{
		Id:                "1",
		EntryPointAddress: "0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789",
		PayMasterAddress:  "0xd93349Ee959d295B115Ee223aF10EF432A8E8523",
		NetWork:           types.Sepolia,
		PayType:           types.PayTypeVerifying,
		EntryPointTag:     types.EntrypointV06,
		Token:             types.ETH,
	}
	MockStrategyMap["2"] = &model.Strategy{
		Id:                "2",
		EntryPointAddress: "0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789",
		PayMasterAddress:  "0xd93349Ee959d295B115Ee223aF10EF432A8E8523",
		PayType:           types.PayTypeERC20,
		EntryPointTag:     types.EntrypointV06,
		NetWork:           types.Sepolia,
		Token:             types.USDT,
	}

	entryPointSupport["0x0576a174D229E3cFA37253523E645A78A0C91B57"] = true
	payMasterSupport["0x0000000000325602a77416A16136FDafd04b299f"] = true
}
func GetStrategyById(strategyId string) *model.Strategy {
	return MockStrategyMap[strategyId]
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
