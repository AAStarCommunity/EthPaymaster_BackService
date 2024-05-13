package model

type SecretConfig struct {
	PriceOracleApiKey string `json:"price_oracle_api_key"`

	NetWorkSecretConfigMap map[string]NetWorkSecretConfig `json:"network_secret_configs"`
}

type NetWorkSecretConfig struct {
	RpcUrl    string `json:"rpc_url"`
	SignerKey string `json:"signer_key"`
}
