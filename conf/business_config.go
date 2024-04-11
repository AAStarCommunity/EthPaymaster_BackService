package conf

import (
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_common/contract/erc20"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"os"
)

var BasicConfig *BusinessConfig
var TokenContractCache map[common.Address]contract_erc20.Contract

func init() {
	originConfig := initBusinessConfig()
	BasicConfig = convertConfig(originConfig)
}
func getBasicConfigPath() *string {
	path := "../conf/business_config.json"
	return &path
}
func convertConfig(originConfig *OriginBusinessConfig) *BusinessConfig {
	basic := &BusinessConfig{}
	basic.NetworkConfigMap = make(map[types.Network]NetWorkConfig)
	basic.SupportEntryPoint = make(map[types.Network]mapset.Set[string])
	basic.SupportPaymaster = make(map[types.Network]mapset.Set[string])
	for network, originNetWorkConfig := range originConfig.NetworkConfigMap {
		//TODO valid
		basic.NetworkConfigMap[network] = NetWorkConfig{
			ChainId:     originNetWorkConfig.ChainId,
			IsTest:      originNetWorkConfig.IsTest,
			RpcUrl:      fmt.Sprintf("%s/%s", originNetWorkConfig.RpcUrl, originNetWorkConfig.ApiKey),
			TokenConfig: originNetWorkConfig.TokenConfig,
		}

		paymasterArr := originConfig.SupportPaymaster[network]
		fmt.Printf("paymasterArr: %v\n", paymasterArr)
		paymasterSet := mapset.NewSet[string]()
		paymasterSet.Append(paymasterArr...)
		basic.SupportPaymaster[network] = paymasterSet

		entryPointArr := originConfig.SupportEntryPoint[network]
		entryPointSet := mapset.NewSet[string]()
		entryPointSet.Append(entryPointArr...)
		basic.SupportEntryPoint[network] = entryPointSet
	}
	return basic
}
func initBusinessConfig() *OriginBusinessConfig {
	var config OriginBusinessConfig
	filePath := getBasicConfigPath()
	file, err := os.Open(*filePath)
	if err != nil {

		panic(fmt.Sprintf("file not found: %s", *filePath))
	}
	//var mapValue map[string]any
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(fmt.Sprintf("parse file error: %s", err))
	}

	return &config
}

type OriginBusinessConfig struct {
	NetworkConfigMap  map[types.Network]*OriginNetWorkConfig `json:"network_config"`
	SupportEntryPoint map[types.Network][]string             `json:"support_entrypoint"`
	SupportPaymaster  map[types.Network][]string             `json:"support_paymaster"`
}
type OriginNetWorkConfig struct {
	ChainId     string                     `json:"chain_id"`
	IsTest      bool                       `json:"is_test"`
	RpcUrl      string                     `json:"rpc_url"`
	ApiKey      string                     `json:"api_key"`
	TokenConfig map[types.TokenType]string `json:"token_config"`
	GasToken    types.TokenType
}

type BusinessConfig struct {
	NetworkConfigMap  map[types.Network]NetWorkConfig      `json:"network_config"`
	SupportEntryPoint map[types.Network]mapset.Set[string] `json:"support_entrypoint"`
	SupportPaymaster  map[types.Network]mapset.Set[string] `json:"support_paymaster"`
}
type NetWorkConfig struct {
	ChainId     string                     `json:"chain_id"`
	IsTest      bool                       `json:"is_test"`
	RpcUrl      string                     `json:"rpc_url"`
	TokenConfig map[types.TokenType]string `json:"token_config"`
	GasToken    types.TokenType
}

func GetTokenAddress(networkParam types.Network, tokenParam types.TokenType) string {
	networkConfig := BasicConfig.NetworkConfigMap[networkParam]

	return networkConfig.TokenConfig[tokenParam]
}
func CheckEntryPointExist(network2 types.Network, address string) bool {
	entryPointSet := BasicConfig.SupportEntryPoint[network2]
	entryPointSetValue := entryPointSet
	return entryPointSetValue.Contains(address)
}
func GetGasToken(networkParam types.Network) types.TokenType {
	networkConfig := BasicConfig.NetworkConfigMap[networkParam]
	return networkConfig.GasToken
}

func GetChainId(newworkParam types.Network) string {
	networkConfig := BasicConfig.NetworkConfigMap[newworkParam]
	return networkConfig.ChainId
}
func GetEthereumRpcUrl(network types.Network) string {
	networkConfig := BasicConfig.NetworkConfigMap[network]
	return networkConfig.RpcUrl
}

var (
	testNetWork = mapset.NewSet(
		types.Sepolia, types.OptimismSepolia, types.ArbitrumSeplia, types.ScrollSepolia, types.StarknetSepolia, types.BaseSepolia)
	opeStackNetWork = mapset.NewSet(
		types.Optimism, types.OptimismSepolia, types.Base, types.BaseSepolia)
)

func IsTestNet(network types.Network) bool {
	return testNetWork.Contains(network)
}
func IsOpStackNetWork(network types.Network) bool {
	return opeStackNetWork.Contains(network)
}
