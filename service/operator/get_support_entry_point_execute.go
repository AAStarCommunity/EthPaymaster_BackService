package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/config"
)

func GetSupportEntrypointExecute(networkStr string) ([]string, error) {

	entryPoints, err := config.GetSupportEntryPoints(global_const.Network(networkStr))
	if err != nil {
		return nil, err
	}

	it := entryPoints.Iterator()
	var entrypointArr []string
	for entry := range it.C {
		entrypointArr = append(entrypointArr, entry)
	}
	return entrypointArr, nil
}
