package dashboard_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"errors"
	"github.com/ethereum/go-ethereum/common"
)

// TODO just Temp Mock
var MockStrategyMap = map[string]*model.Strategy{}
var payMasterSupport = map[string]bool{}
var entryPointSupport = map[string]bool{}

func init() {
	entrypoint := common.HexToAddress("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789")
	paymaster := common.HexToAddress("0xAEbF4C90b571e7D5cb949790C9b8Dc0280298b63")
	MockStrategyMap["1"] = &model.Strategy{
		Id: "1",
		NetWorkInfo: &model.NetWorkInfo{
			NetWork: types.Sepolia,
			Token:   types.ETH,
		},
		EntryPointInfo: &model.EntryPointInfo{
			EntryPointAddress: &entrypoint,
			EntryPointTag:     types.EntrypointV06,
		},
		ExecuteRestriction: model.StrategyExecuteRestriction{
			EffectiveStartTime: 1710044496,
			EffectiveEndTime:   1820044496,
		},
		PaymasterInfo: &model.PaymasterInfo{
			PayMasterAddress: &paymaster,
			PayType:          types.PayTypeVerifying,
		},
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
