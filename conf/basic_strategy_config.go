package conf

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"os"
	"strings"
)

var originConfig *OriginBusinessConfig
var BasicStrategyConfig map[string]*model.Strategy = make(map[string]*model.Strategy)
var SuitableStrategyMap map[types.Network]map[string]map[types.PayType]*model.Strategy = make(map[types.Network]map[string]map[types.PayType]*model.Strategy)

func getStrategyConfigPath() *string {
	path := fmt.Sprintf("../conf/basic_strategy_%s_config.json", strings.ToLower(Environment.Name))
	fmt.Println(path)
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		path = fmt.Sprintf("../conf/basic_strategy_config.json")
	}
	return &path
}
func BasicStrategyConfigInit() {
	filePah := getStrategyConfigPath()
	file, err := os.Open(*filePah)
	if err != nil {
		panic(fmt.Sprintf("file not found: %s", *filePah))
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
	BasicStrategyConfig = strateyMap
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
		SuitableStrategyMap[strategy.NetWorkInfo.NetWork][strategy.GetEntryPointAddress().String()][strategy.GetPayType()] = strategy
	}
	return config, nil
}
