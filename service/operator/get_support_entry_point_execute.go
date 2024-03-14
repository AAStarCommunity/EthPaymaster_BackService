package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
)

func GetSupportEntrypointExecute(network string) (*[]model.EntrypointDomain, error) {
	entrypoints := make([]model.EntrypointDomain, 0)
	entrypoints = append(entrypoints, model.EntrypointDomain{
		Address:    "0x0576a174D229E3cFA37253523E645A78A0C91B57",
		Desc:       "desc",
		NetWork:    types.Sepolia,
		StrategyId: "1",
	})
	return &entrypoints, nil
}
