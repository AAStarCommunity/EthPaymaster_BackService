package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
)

func GetSupportEntrypointExecute(networkStr string) (*[]model.EntrypointDomain, error) {
	entrypoints := make([]model.EntrypointDomain, 0)
	entrypoints = append(entrypoints, model.EntrypointDomain{
		Address:    "0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789",
		Desc:       "desc",
		NetWork:    global_const.EthereumSepolia,
		StrategyId: "1",
	})
	return &entrypoints, nil
}
