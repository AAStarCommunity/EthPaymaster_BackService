package operator

import "AAStarCommunity/EthPaymaster_BackService/common/model"

func GetSupportEntrypointExecute(network string) (*model.GetSupportEntryPointResponse, error) {
	return &model.GetSupportEntryPointResponse{
		EntrypointDomains: make([]model.EntrypointDomain, 0),
	}, nil
}
