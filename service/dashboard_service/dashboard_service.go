package dashboard_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"encoding/json"
	"errors"
	mapset "github.com/deckarep/golang-set/v2"
	"golang.org/x/time/rate"
	"golang.org/x/xerrors"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math/big"
	"sync"
)

var (
	dashBoardDb *gorm.DB
	onlyOnce    = sync.Once{}
)

func Init() {
	onlyOnce.Do(func() {
		configDBDsn := config.GetConfigDBDSN()

		configDBVar, err := gorm.Open(postgres.Open(configDBDsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		dashBoardDb = configDBVar
	})

}

type StrategyDBModel struct {
	model.BaseData
	DeletedAt          gorm.DeletedAt              `gorm:"softDelete:flag" json:"deleted_at"`
	Description        string                      `gorm:"type:varchar(500)" json:"description"`
	StrategyCode       string                      `gorm:"type:varchar(255)" json:"strategy_code"`
	ProjectCode        string                      `gorm:"type:varchar(255)" json:"project_code"`
	StrategyName       string                      `gorm:"type:varchar(255)" json:"strategy_name"`
	UserId             string                      `gorm:"type:varchar(255)" json:"user_id"`
	Status             global_const.StrategyStatus `gorm:"type:varchar(20)" json:"status"`
	ExecuteRestriction datatypes.JSON              `gorm:"type:json" json:"execute_restriction"`
	Extra              datatypes.JSON              `gorm:"type:json" json:"extra"`
}

func (StrategyDBModel) TableName() string {
	return config.GetStrategyConfigTableName()
}

// GetStrategyByCode is Sponsor Type , need GasTank
func GetStrategyByCode(strategyCode string, entryPointVersion global_const.EntrypointVersion, network global_const.Network) (*model.Strategy, error) {
	if entryPointVersion == "" {
		entryPointVersion = global_const.EntrypointV06
	}
	strategyDbModel := &StrategyDBModel{}
	tx := dashBoardDb.Where("strategy_code = ?", strategyCode).First(&strategyDbModel)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, xerrors.Errorf("strategy not found: %w", tx.Error)
		} else {
			return nil, xerrors.Errorf("error when finding strategy: %w", tx.Error)
		}
	}
	strategy, err := convertStrategyDBModelToStrategy(strategyDbModel, entryPointVersion, network)
	if err != nil {
		return nil, err
	}

	strategy.ProjectSponsor = true

	return strategy, nil
}

func convertStrategyDBModelToStrategy(strategyDBModel *StrategyDBModel, entryPointVersion global_const.EntrypointVersion, network global_const.Network) (*model.Strategy, error) {
	entryPointAddress := config.GetEntrypointAddress(network, entryPointVersion)

	if entryPointAddress == nil {
		return nil, errors.New("entryPointAddress not found")
	}
	paymasterAddress := config.GetPaymasterAddress(network, entryPointVersion)
	if paymasterAddress == nil {
		return nil, errors.New("paymasterAddress not found")
	}

	if strategyDBModel.Status == "" {
		strategyDBModel.Status = global_const.StrategyStatusDisable
	}
	strategyExecuteRestrictionJson := StrategyExecuteRestrictionJson{}
	if strategyDBModel.ExecuteRestriction != nil {
		eJson, _ := strategyDBModel.ExecuteRestriction.MarshalJSON()
		err := json.Unmarshal(eJson, &strategyExecuteRestrictionJson)
		if err != nil {
			return nil, xerrors.Errorf("error when unmarshal strategyExecuteRestriction: %w", err)
		}

		if err != nil {
			return nil, xerrors.Errorf("error when unmarshal strategyExecuteRestriction: %w", err)
		}
	}
	strategyExecuteRestriction := &model.StrategyExecuteRestriction{
		EffectiveStartTime: big.NewInt(strategyExecuteRestrictionJson.EffectiveStartTime),
		EffectiveEndTime:   big.NewInt(strategyExecuteRestrictionJson.EffectiveEndTime),
		GlobalMaxUSD:       big.NewFloat(strategyExecuteRestrictionJson.GlobalMaxUSD),
		GlobalMaxOpCount:   big.NewInt(strategyExecuteRestrictionJson.GlobalMaxOpCount),
		DayMaxUSD:          big.NewFloat(strategyExecuteRestrictionJson.DayMaxUSD),
		Status:             strategyDBModel.Status,
	}
	if strategyExecuteRestrictionJson.BanSenderAddress != nil {
		strategyExecuteRestriction.BanSenderAddress = mapset.NewSetWithSize[string](len(strategyExecuteRestrictionJson.BanSenderAddress))
		for _, v := range strategyExecuteRestrictionJson.BanSenderAddress {
			strategyExecuteRestriction.BanSenderAddress.Add(v)
		}
	}
	if strategyExecuteRestrictionJson.AccessProject != nil {
		strategyExecuteRestriction.AccessProject = mapset.NewSetWithSize[string](len(strategyExecuteRestrictionJson.AccessProject))
		for _, v := range strategyExecuteRestrictionJson.AccessProject {
			strategyExecuteRestriction.AccessProject.Add(v)
		}
	}
	if strategyExecuteRestrictionJson.ChainIdWhiteList != nil {
		strategyExecuteRestriction.ChainIdWhiteList = mapset.NewSetWithSize[string](len(strategyExecuteRestrictionJson.ChainIdWhiteList))
		for _, v := range strategyExecuteRestrictionJson.ChainIdWhiteList {
			strategyExecuteRestriction.ChainIdWhiteList.Add(v)
		}
	}

	return &model.Strategy{
		StrategyCode: strategyDBModel.StrategyCode,
		Description:  strategyDBModel.Description,
		NetWorkInfo: &model.NetWorkInfo{
			NetWork:  network,
			GasToken: config.GetGasToken(network),
		},
		EntryPointInfo: &model.EntryPointInfo{
			EntryPointVersion: entryPointVersion,
			EntryPointAddress: config.GetEntrypointAddress(network, entryPointVersion),
		},
		PaymasterInfo: &model.PaymasterInfo{
			PayMasterAddress:        config.GetPaymasterAddress(network, entryPointVersion),
			PayType:                 global_const.PayTypeProjectSponsor,
			IsProjectErc20PayEnable: false,
		},
		ExecuteRestriction: strategyExecuteRestriction,
	}, nil
}

type StrategyExecuteRestrictionJson struct {
	BanSenderAddress   []string `json:"ban_sender_address"`
	EffectiveStartTime int64    `json:"start_time"`
	EffectiveEndTime   int64    `json:"end_time"`
	GlobalMaxUSD       float64  `json:"global_max_usd"`
	GlobalMaxOpCount   int64    `json:"global_max_op_count"`
	DayMaxUSD          float64  `json:"day_max_usd"`
	AccessProject      []string `json:"access_project"`
	AccessErc20        []string `json:"access_erc20"`
	ChainIdWhiteList   []string `json:"chain_id_whitelist"`
}

// GetSuitableStrategyWithOutCode get suitable strategy by entryPointVersion, chain,
//
//	For Offical StrategyConfig,
func GetSuitableStrategyWithOutCode(entryPointVersion global_const.EntrypointVersion, chain global_const.Network, gasUseToken global_const.TokenType) (*model.Strategy, error) {
	if entryPointVersion == "" {
		entryPointVersion = global_const.EntrypointV06
	}
	gasToken := config.GetGasToken(chain)
	entryPointAddress := config.GetEntrypointAddress(chain, entryPointVersion)
	paymasterAddress := config.GetPaymasterAddress(chain, entryPointVersion)
	payType := global_const.PayTypeUserSponsor
	isPerc20Enable := false
	if gasUseToken != "" {
		payType = global_const.PayTypeERC20
		if config.IsPErc20Token(gasUseToken) {
			isPerc20Enable = true
		}
	}

	strategy := &model.Strategy{
		NetWorkInfo: &model.NetWorkInfo{
			NetWork:  chain,
			GasToken: gasToken,
		},
		EntryPointInfo: &model.EntryPointInfo{
			EntryPointVersion: entryPointVersion,
			EntryPointAddress: entryPointAddress,
		},
		PaymasterInfo: &model.PaymasterInfo{
			PayMasterAddress:        paymasterAddress,
			PayType:                 payType,
			IsProjectErc20PayEnable: isPerc20Enable,
		},
		ExecuteRestriction: &model.StrategyExecuteRestriction{
			Status: global_const.StrategyStatusAchieve,
		},
		Erc20TokenType: gasUseToken,
	}
	return strategy, nil
}

func IsEntryPointsSupport(address string, chain global_const.Network) bool {
	supportEntryPointSet, _ := config.GetSupportEntryPoints(chain)
	if supportEntryPointSet == nil {
		return false
	}
	return supportEntryPointSet.Contains(address)
}
func IsPayMasterSupport(address string, chain global_const.Network) bool {
	supportPayMasterSet, _ := config.GetSupportPaymaster(chain)
	if supportPayMasterSet == nil {
		return false
	}

	return supportPayMasterSet.Contains(address)
}

type ApiKeyDbModel struct {
	model.BaseData
	UserId    int64          `gorm:"column:user_id;type:integer" json:"user_id"`
	Disable   bool           `gorm:"column:disable;type:bool" json:"disable"`
	ApiKey    string         `gorm:"column:api_key;type:varchar(255)" json:"api_key"`
	KeyName   string         `gorm:"column:key_name;type:varchar(255)" json:"key_name"`
	DeletedAt gorm.DeletedAt `gorm:"softDelete:flag" json:"deleted_at"`
	Extra     datatypes.JSON `gorm:"column:extra" json:"extra"`
}

func (*ApiKeyDbModel) TableName() string {
	return config.GetAPIKeyTableName()
}

func (m *ApiKeyDbModel) GetRateLimit() rate.Limit {
	return 10
}
func convertApiKeyDbModelToApiKeyModel(apiKeyDbModel *ApiKeyDbModel) *model.ApiKeyModel {
	return &model.ApiKeyModel{
		Disable:   apiKeyDbModel.Disable,
		ApiKey:    apiKeyDbModel.ApiKey,
		RateLimit: 10,
		UserId:    apiKeyDbModel.UserId,
	}
}
func GetAPiInfoByApiKey(apiKey string) (*model.ApiKeyModel, error) {
	apikeyModel := &ApiKeyDbModel{}
	tx := dashBoardDb.Where("api_key = ?", apiKey).First(&apikeyModel)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, tx.Error
		}
		return nil, xerrors.Errorf("error when finding apikey: %w", tx.Error)
	}
	apikeyRes := convertApiKeyDbModelToApiKeyModel(apikeyModel)
	return apikeyRes, nil
}

type PaymasterRecallLogDbModel struct {
	model.BaseData
	ProjectUserId   int64          `gorm:"column:project_user_id;type:integer" json:"project_user_id"`
	ProjectApikey   string         `gorm:"column:project_apikey;type:varchar(255)" json:"project_apikey"`
	PaymasterMethod string         `gorm:"column:paymaster_method;type:varchar(25)" json:"paymaster_method"`
	SendTime        string         `gorm:"column:send_time;type:varchar(50)" json:"send_time"`
	Latency         int64          `gorm:"column:latency;type:integer" json:"latency"`
	RequestBody     string         `gorm:"column:request_body;type:varchar(500)" json:"request_body"`
	ResponseBody    string         `gorm:"column:response_body;type:varchar(1000)" json:"response_body"`
	NetWork         string         `gorm:"column:network;type:varchar(25)" json:"network"`
	Status          int            `gorm:"column:status;type:integer" json:"status"`
	Extra           datatypes.JSON `gorm:"column:extra" json:"extra"`
}

func (*PaymasterRecallLogDbModel) TableName() string {
	return "paymaster_recall_log"
}
func CreatePaymasterCall(recallModel *PaymasterRecallLogDbModel) error {
	return dashBoardDb.Create(recallModel).Error
}
