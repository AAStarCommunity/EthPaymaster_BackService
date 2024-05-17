package config

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"encoding/json"
	"fmt"
	"golang.org/x/xerrors"
	"math/big"
	"os"
)

var basicStrategyConfig = make(map[string]*model.Strategy)

// suitableStrategyMap[chain][entrypoint][payType]
var suitableStrategyMap = make(map[global_const.Network]map[global_const.EntrypointVersion]map[global_const.PayType]*model.Strategy)

func GetBasicStrategyConfig(strategyCode global_const.BasicStrategyCode) *model.Strategy {
	strategy := basicStrategyConfig[string(strategyCode)]
	paymasterAddress := GetPaymasterAddress(strategy.GetNewWork(), strategy.GetStrategyEntrypointVersion())
	strategy.PaymasterInfo.PayMasterAddress = &paymasterAddress
	entryPointAddress := GetEntrypointAddress(strategy.GetNewWork(), strategy.GetStrategyEntrypointVersion())
	strategy.EntryPointInfo.EntryPointAddress = &entryPointAddress
	return strategy

}
func GetSuitableStrategy(entrypointVersion global_const.EntrypointVersion, chain global_const.Network, payType global_const.PayType) (*model.Strategy, error) {
	//TODO
	strategy := suitableStrategyMap[chain][entrypointVersion][payType]
	if strategy == nil {
		return nil, xerrors.Errorf("strategy not found")
	}
	return strategy, nil
}

func basicStrategyInit(path string) {
	if path == "" {
		panic("pathParam is empty")
	}

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	//var mapValue map[string]any
	decoder := json.NewDecoder(file)
	var config map[string]map[string]any
	err = decoder.Decode(&config)
	if err != nil {
		panic(fmt.Sprintf("parse file error: %s", err))
	}
	strateyMap, err := convertMapToStrategyConfig(config)
	if err != nil {
		panic(fmt.Sprintf("parse file error: %s", err))
	}
	basicStrategyConfig = strateyMap
}
func convertMapToStrategyConfig(data map[string]map[string]any) (map[string]*model.Strategy, error) {
	config := make(map[string]*model.Strategy)

	for key, value := range data {
		effectiveStartTime, ok := new(big.Int).SetString(value["effective_start_time"].(string), 10)
		if !ok {
			return nil, xerrors.Errorf("effective_start_time illegal")
		}
		effectiveEndTime, ok := new(big.Int).SetString(value["effective_end_time"].(string), 10)
		if !ok {
			return nil, xerrors.Errorf("effective_end_time illegal")
		}
		accessProjectStr := value["access_project"].(string)
		strategy := &model.Strategy{
			Id:           key,
			StrategyCode: key,
			NetWorkInfo: &model.NetWorkInfo{
				NetWork: global_const.Network(value["network"].(string)),
			},
			EntryPointInfo: &model.EntryPointInfo{
				EntryPointVersion: global_const.EntrypointVersion(value["entrypoint_version"].(string)),
			},

			ExecuteRestriction: model.StrategyExecuteRestriction{
				EffectiveStartTime: effectiveStartTime,
				EffectiveEndTime:   effectiveEndTime,
				AccessProject:      utils.ConvertStringToSet(accessProjectStr, ","),
			},
			PaymasterInfo: &model.PaymasterInfo{
				PayType: global_const.PayType(value["paymaster_pay_type"].(string)),
			},
		}
		if strategy.GetPayType() == global_const.PayTypeERC20 {
			erc20TokenStr := value["access_erc20"].(string)
			strategy.NetWorkInfo.GasToken = global_const.TokenType(erc20TokenStr)
			strategy.ExecuteRestriction.AccessErc20 = utils.ConvertStringToSet(erc20TokenStr, ",")
		}

		config[key] = strategy
		if suitableStrategyMap[strategy.NetWorkInfo.NetWork] == nil {
			suitableStrategyMap[strategy.NetWorkInfo.NetWork] = make(map[global_const.EntrypointVersion]map[global_const.PayType]*model.Strategy)
		}
		if suitableStrategyMap[strategy.NetWorkInfo.NetWork][strategy.GetStrategyEntrypointVersion()] == nil {
			suitableStrategyMap[strategy.NetWorkInfo.NetWork][strategy.GetStrategyEntrypointVersion()] = make(map[global_const.PayType]*model.Strategy)
		}
		suitableStrategyMap[strategy.NetWorkInfo.NetWork][strategy.GetStrategyEntrypointVersion()][strategy.GetPayType()] = strategy
	}
	return config, nil
}
