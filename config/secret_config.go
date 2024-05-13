package config

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"encoding/json"
	"fmt"
	"os"
)

var secretConfig *model.SecretConfig
var signerConfig = make(SignerConfigMap)

type SignerConfigMap map[global_const.Network]*global_const.EOA

func secretConfigInit(secretConfigPath string) {
	if secretConfigPath == "" {
		panic("secretConfigPath is empty")
	}
	file, err := os.Open(secretConfigPath)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	var config model.SecretConfig
	err = decoder.Decode(&config)
	if err != nil {
		panic(fmt.Sprintf("parse file error: %s", err))
	}
	secretConfig = &config
	for network, originNetWorkConfig := range secretConfig.NetWorkSecretConfigMap {
		signerKey := originNetWorkConfig.SignerKey
		eoa, err := global_const.NewEoa(signerKey)
		if err != nil {
			panic(fmt.Sprintf("signer key error: %s", err))
		}

		signerConfig[global_const.Network(network)] = eoa
	}
}
func GetNetworkSecretConfig(network global_const.Network) model.NetWorkSecretConfig {
	return secretConfig.NetWorkSecretConfigMap[string(network)]
}

func GetPriceOracleApiKey() string {
	return secretConfig.PriceOracleApiKey
}
func GetNewWorkClientURl(network global_const.Network) string {
	return secretConfig.NetWorkSecretConfigMap[string(network)].RpcUrl
}
func GetSignerKey(network global_const.Network) string {
	return secretConfig.NetWorkSecretConfigMap[string(network)].SignerKey
}
func GetSigner(network global_const.Network) *global_const.EOA {
	return signerConfig[network]
}
