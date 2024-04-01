package dashboard_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/erc20_token"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"errors"
	"github.com/ethereum/go-ethereum/common"
)

// TODO just Temp Mock
var MockStrategyMap = map[string]*model.Strategy{}
var payMasterSupport = map[string]bool{}
var entryPointSupport = map[string]bool{}

func init() {
	MockStrategyMap["1"] = &model.Strategy{
		Id: "1",
		NetWorkInfo: &model.NetWorkInfo{
			NetWork: network.Sepolia,
			Token:   erc20_token.ETH,
		},
		EntryPointInfo: &model.EntryPointInfo{
			EntryPointAddress: common.HexToAddress("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789"),
			EntryPointTag:     types.EntrypointV06,
		},
		ExecuteRestriction: model.StrategyExecuteRestriction{
			StartTime: 1710044496,
			EndTime:   1820044496,
		},
		PaymasterInfo: &model.PaymasterInfo{
			PayMasterAddress: common.HexToAddress("0xAEbF4C90b571e7D5cb949790C9b8Dc0280298b63"),
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
func GetSuitableStrategy(entrypoint string, chain network.Network, token string) (*model.Strategy, error) {
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
