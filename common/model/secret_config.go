package model

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
	TimeZone string `json:"tz"`
	SslMode  string `json:"ssl_mode"`
}

type SecretConfig struct {
	PriceOracleApiKey string `json:"price_oracle_api_key"`

	NetWorkSecretConfigMap map[string]NetWorkSecretConfig `json:"network_secret_configs"`

	ConfigDBConfig          DBConfig      `json:"config_db_config"`
	RelayDBConfig           DBConfig      `json:"relay_db_config"`
	ApiKeyTableName         string        `json:"api_key_table_name"`
	StrategyConfigTableName string        `json:"strategy_config_table_name"`
	SponsorConfig           SponsorConfig `json:"sponsor_config"`
}
type SponsorConfig struct {
	SponsorDepositAddress     string   `json:"sponsor_deposit_address"`
	SponsorWithdrawPrivateKey string   `json:"sponsor_withdraw_private_key"`
	DashBoardSignerAddress    string   `json:"dashboard_signer_address"`
	SponsorTestClientUrl      string   `json:"sponsor_client_rpc_test_net"`
	SponsorMainClientUrl      string   `json:"sponsor_client_rpc_main_net"`
	FreeSponsorWhitelist      []string `json:"free_sponsor_whitelist"`
	FreeSponsorAPIList        []string `json:"free_sponsor_white_api_list"`
}

type NetWorkSecretConfig struct {
	RpcUrl    string `json:"rpc_url"`
	SignerKey string `json:"signer_key"`
}
