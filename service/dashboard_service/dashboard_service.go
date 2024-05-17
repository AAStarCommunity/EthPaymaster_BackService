package dashboard_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"errors"
	"golang.org/x/time/rate"
	"golang.org/x/xerrors"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

var (
	configDB *gorm.DB
	relayDB  *gorm.DB
	onlyOnce = sync.Once{}
)

func Init() {
	onlyOnce.Do(func() {
		configDBDsn := config.GetConfigDBDSN()
		relayDBDsn := config.GetRelayDBDSN()

		configDBVar, err := gorm.Open(postgres.Open(configDBDsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		configDB = configDBVar

		relayDBVar, err := gorm.Open(postgres.Open(relayDBDsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		relayDB = relayDBVar
	})

}

type StrategyDBModel struct {
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
func GetStrategyByCode(strategyCode string, entryPointVersion global_const.EntrypointVersion, chain global_const.Network) (*model.Strategy, error) {

	strategyDbModel := &StrategyDBModel{}
	tx := configDB.Where("strategy_code = ?", strategyCode).First(&strategyDbModel)
	if tx.Error != nil {
		return nil, tx.Error
	}

	strategy, err := convertStrategyDBModelToStrategy(strategyDbModel)
	if err != nil {
		return nil, err
	}
	paymasterAddress := config.GetPaymasterAddress(strategy.GetNewWork(), strategy.GetStrategyEntrypointVersion())
	strategy.PaymasterInfo.PayMasterAddress = &paymasterAddress
	entryPointAddress := config.GetEntrypointAddress(strategy.GetNewWork(), strategy.GetStrategyEntrypointVersion())
	strategy.EntryPointInfo.EntryPointAddress = &entryPointAddress
	return strategy, nil
}

func convertStrategyDBModelToStrategy(strategyDBModel *StrategyDBModel) (*model.Strategy, error) {
	return &model.Strategy{}, nil
}

// GetSuitableStrategy get suitable strategy by entryPointVersion, chain,
//
//	For Offical StrategyConfig,
func GetSuitableStrategy(entryPointVersion global_const.EntrypointVersion, chain global_const.Network, gasUseToken global_const.TokenType) (*model.Strategy, error) {
	if entryPointVersion == "" {
		entryPointVersion = global_const.EntrypointV06
	}
	gasToken := config.GetGasToken(chain)
	entryPointAddress := config.GetEntrypointAddress(chain, entryPointVersion)
	paymasterAddress := config.GetPaymasterAddress(chain, entryPointVersion)
	payType := global_const.PayTypeVerifying
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
			EntryPointAddress: &entryPointAddress,
		},
		PaymasterInfo: &model.PaymasterInfo{
			PayMasterAddress:        &paymasterAddress,
			PayType:                 payType,
			IsProjectErc20PayEnable: isPerc20Enable,
		},
		Erc20TokenType: gasUseToken,
	}
	if strategy == nil {
		return nil, errors.New("strategy not found")
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
	Disable bool           `gorm:"column:disable;type:bool" json:"disable"`
	ApiKey  string         `gorm:"column:api_key;type:varchar(255)" json:"api_key"`
	KeyName string         `gorm:"column:key_name;type:varchar(255)" json:"key_name"`
	Extra   datatypes.JSON `gorm:"column:extra" json:"extra"`
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
	}
}
func GetAPiInfoByApiKey(apiKey string) (*model.ApiKeyModel, error) {
	apikeyModel := &ApiKeyDbModel{}
	tx := configDB.Where("api_key = ?", apiKey).First(&apikeyModel)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, tx.Error
		}
		return nil, xerrors.Errorf("error when finding apikey: %w", tx.Error)
	}
	apikeyRes := convertApiKeyDbModelToApiKeyModel(apikeyModel)
	return apikeyRes, nil
}
