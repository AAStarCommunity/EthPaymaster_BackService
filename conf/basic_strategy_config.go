package conf

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/xerrors"
	"math/big"
	"os"
)

var basicStrategyConfig = make(map[string]*model.Strategy)

// suitableStrategyMap[chain][entrypoint][payType]
var suitableStrategyMap = make(map[types.Network]map[string]map[types.PayType]*model.Strategy)

func GetBasicStrategyConfig(key string) *model.Strategy {
	return basicStrategyConfig[key]
}
func GetSuitableStrategy(entrypoint string, chain types.Network, payType types.PayType) (*model.Strategy, error) {
	strategy := suitableStrategyMap[chain][entrypoint][payType]
	if strategy == nil {
		return nil, xerrors.Errorf("strategy not found")
	}
	return strategy, nil
}

func BasicStrategyInit(path string) {
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
		paymasterAddress := common.HexToAddress(value["paymaster_address"].(string))
		entryPointAddress := common.HexToAddress(value["entrypoint_address"].(string))
		effectiveStartTime, ok := new(big.Int).SetString(value["effective_start_time"].(string), 10)
		if !ok {
			return nil, xerrors.Errorf("effective_start_time illegal")
		}
		effectiveEndTime, ok := new(big.Int).SetString(value["effective_end_time"].(string), 10)
		if !ok {
			return nil, xerrors.Errorf("effective_end_time illegal")
		}

		strategy := &model.Strategy{
			Id:           key,
			StrategyCode: key,
			NetWorkInfo: &model.NetWorkInfo{
				NetWork: types.Network(value["network"].(string)),
			},
			EntryPointInfo: &model.EntryPointInfo{
				EntryPointAddress: &entryPointAddress,
				EntryPointTag:     types.EntrypointVersion(value["entrypoint_tag"].(string)),
			},

			ExecuteRestriction: model.StrategyExecuteRestriction{
				EffectiveStartTime: effectiveStartTime,
				EffectiveEndTime:   effectiveEndTime,
			},
			PaymasterInfo: &model.PaymasterInfo{
				PayMasterAddress: &paymasterAddress,
				PayType:          types.PayType(value["paymaster_pay_type"].(string)),
			},
		}

		config[key] = strategy
		if suitableStrategyMap[strategy.NetWorkInfo.NetWork] == nil {
			suitableStrategyMap[strategy.NetWorkInfo.NetWork] = make(map[string]map[types.PayType]*model.Strategy)
		}
		if suitableStrategyMap[strategy.NetWorkInfo.NetWork][strategy.GetEntryPointAddress().String()] == nil {
			suitableStrategyMap[strategy.NetWorkInfo.NetWork][strategy.GetEntryPointAddress().String()] = make(map[types.PayType]*model.Strategy)
		}
		suitableStrategyMap[strategy.NetWorkInfo.NetWork][strategy.GetEntryPointAddress().String()][strategy.GetPayType()] = strategy
	}
	return config, nil
}
