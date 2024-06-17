package config

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"context"
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"math/big"
	"os"
	"sync"
)

var dsnTemplate = "host=%s port=%v user=%s password=%s dbname=%s TimeZone=%s sslmode=%s"

var secretConfig *model.SecretConfig
var signerConfig = make(SignerConfigMap)
var depositer *global_const.EOA

var sponsorTestNetClient *ethclient.Client
var sponsorTestNetClientChainId *big.Int
var sponsorMainNetClient *ethclient.Client
var sponsorMainNetClientChainId *big.Int
var onlyOnce = sync.Once{}

func GetPaymasterSponsorClient(isTestNet bool) *ethclient.Client {
	if isTestNet {
		return sponsorTestNetClient

	}
	return sponsorMainNetClient
}
func GetPaymasterSponsorChainId(isTestNet bool) *big.Int {
	if isTestNet {
		return sponsorTestNetClientChainId
	}
	return sponsorMainNetClientChainId
}

var sponsorWhitelist = mapset.NewSet[string]()

type SignerConfigMap map[global_const.Network]*global_const.EOA

func GetDepositer() *global_const.EOA {
	return depositer
}
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
	depositer, err = global_const.NewEoa(secretConfig.SponsorConfig.SponsorDepositPrivateKey)
	if err != nil {
		panic(fmt.Sprintf("signer key error: %s", err))
	}
	onlyOnce.Do(func() {
		paymasterSponsorMainNetClient, err := ethclient.Dial(secretConfig.SponsorConfig.SponsorMainClientUrl)
		if err != nil {
			panic(fmt.Sprintf("paymaster inner client error: %s", err))
		}
		paymasterInnerClientChainId, err := paymasterSponsorMainNetClient.ChainID(context.Background())
		if err != nil {
			panic(fmt.Sprintf("paymaster inner client chain id error: %s", err))
		}
		sponsorMainNetClient = paymasterSponsorMainNetClient
		sponsorMainNetClientChainId = paymasterInnerClientChainId

		paymasterSponsorTestNetClient, err := ethclient.Dial(secretConfig.SponsorConfig.SponsorTestClientUrl)
		if err != nil {
			panic(fmt.Sprintf("paymaster inner client error: %s", err))
		}
		paymasterInnerClientChainId, err = paymasterSponsorTestNetClient.ChainID(context.Background())
		if err != nil {
			panic(fmt.Sprintf("paymaster inner client chain id error: %s", err))
		}
		sponsorTestNetClient = paymasterSponsorTestNetClient
		sponsorTestNetClientChainId = paymasterInnerClientChainId
	})
	logrus.Debugf("secretConfig [%v]", secretConfig)
	if secretConfig.SponsorConfig.FreeSponsorWhitelist != nil {
		sponsorWhitelist.Append(secretConfig.SponsorConfig.FreeSponsorWhitelist...)
	}
}

func IsSponsorWhitelist(senderAddress string) bool {
	return sponsorWhitelist.Contains(senderAddress)
}
func GetNetworkSecretConfig(network global_const.Network) model.NetWorkSecretConfig {
	return secretConfig.NetWorkSecretConfigMap[string(network)]
}

func CheckNetworkSupport(network global_const.Network) bool {
	_, ok := secretConfig.NetWorkSecretConfigMap[string(network)]
	return ok
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
func GetAPIKeyTableName() string {
	return secretConfig.ApiKeyTableName
}
func GetSponsorConfig() *model.SponsorConfig {
	//TODO
	return &secretConfig.SponsorConfig
}
func GetStrategyConfigTableName() string {
	return secretConfig.StrategyConfigTableName
}
func GetConfigDBDSN() string {
	return fmt.Sprintf(dsnTemplate,
		secretConfig.ConfigDBConfig.Host,
		secretConfig.ConfigDBConfig.Port,
		secretConfig.ConfigDBConfig.User,
		secretConfig.ConfigDBConfig.Password,
		secretConfig.ConfigDBConfig.DBName,
		secretConfig.ConfigDBConfig.TimeZone,
		secretConfig.ConfigDBConfig.SslMode)
}
func GetRelayDBDSN() string {
	return fmt.Sprintf(dsnTemplate,
		secretConfig.RelayDBConfig.Host,
		secretConfig.RelayDBConfig.Port,
		secretConfig.RelayDBConfig.User,
		secretConfig.RelayDBConfig.Password,
		secretConfig.RelayDBConfig.DBName,
		secretConfig.RelayDBConfig.TimeZone,
		secretConfig.RelayDBConfig.SslMode)
}
