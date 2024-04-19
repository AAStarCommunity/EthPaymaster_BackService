package conf

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/envirment"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/xerrors"
	"os"
	"strings"
	"sync"
)

var once sync.Once
var basicStrategyConfig map[string]*model.Strategy
var suitableStrategyMap map[types.Network]map[string]map[types.PayType]*model.Strategy

func GetBasicStrategyConfig(key string) *model.Strategy {
	once.Do(func() {
		if basicStrategyConfig == nil {
			BasicStrategyInit()
		}
	})
	return basicStrategyConfig[key]
}
func GetSuitableStrategy(entrypoint string, chain types.Network, payType types.PayType) (*model.Strategy, error) {
	once.Do(func() {
		if basicStrategyConfig == nil {
			BasicStrategyInit()
		}
	})
	strategy := suitableStrategyMap[chain][entrypoint][payType]
	if strategy == nil {
		return nil, xerrors.Errorf("strategy not found")
	}
	return strategy, nil
}

func BasicStrategyInit() {
	basicStrategyConfig = make(map[string]*model.Strategy)
	suitableStrategyMap = make(map[types.Network]map[string]map[types.PayType]*model.Strategy)
	path := fmt.Sprintf("../conf/basic_strategy_%s_config.json", strings.ToLower(envirment.Environment.Name))
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
				EffectiveStartTime: value["effective_start_time"].(int64),
				EffectiveEndTime:   value["effective_end_time"].(int64),
			},
			PaymasterInfo: &model.PaymasterInfo{
				PayMasterAddress: &paymasterAddress,
				PayType:          types.PayType(value["paymaster_pay_type"].(string)),
			},
		}
		config[key] = strategy
		suitableStrategyMap[strategy.NetWorkInfo.NetWork][strategy.GetEntryPointAddress().String()][strategy.GetPayType()] = strategy
	}
	return config, nil
}
