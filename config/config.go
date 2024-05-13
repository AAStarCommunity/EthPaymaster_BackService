package config

func InitConfig(basicStrategyPath string, basicConfigPath string, secretconfig string) {
	secretConfigInit(secretconfig)
	basicStrategyInit(basicStrategyPath)
	basicConfigInit(basicConfigPath)
}
