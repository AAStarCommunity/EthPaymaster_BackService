package conf

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"os"
)

var BasicConfig *BusinessConfig

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
	ChainId          string                     `json:"chain_id"`
	IsTest           bool                       `json:"is_test"`
	RpcUrl           string                     `json:"rpc_url"`
	ApiKey           string                     `json:"api_key"`
	TokenConfig      map[types.TokenType]string `json:"token_config"`
	GasToken         types.TokenType
	GasOracleAddress string
}

type BusinessConfig struct {
	NetworkConfigMap  map[types.Network]NetWorkConfig      `json:"network_config"`
	SupportEntryPoint map[types.Network]mapset.Set[string] `json:"support_entrypoint"`
	SupportPaymaster  map[types.Network]mapset.Set[string] `json:"support_paymaster"`
}
type NetWorkConfig struct {
	ChainId          string                     `json:"chain_id"`
	IsTest           bool                       `json:"is_test"`
	RpcUrl           string                     `json:"rpc_url"`
	TokenConfig      map[types.TokenType]string `json:"token_config"`
	GasToken         types.TokenType
	GasOracleAddress common.Address
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
	TestNetWork = mapset.NewSet(
		types.ETHEREUM_SEPOLIA, types.OPTIMISM_SEPOLIA, types.ARBITRUM_SPEOLIA, types.SCROLL_SEPOLIA, types.STARKET_SEPOLIA, types.BaseSepolia)
	OpeStackNetWork = mapset.NewSet(
		types.OPTIMISM_MAINNET, types.OPTIMISM_SEPOLIA, types.Base, types.BaseSepolia)
	EthereumAdaptableNetWork = mapset.NewSet(
		types.OPTIMISM_MAINNET, types.OPTIMISM_SEPOLIA, types.ETHEREUM_SEPOLIA)
	ArbStackNetWork = mapset.NewSet(
		types.ARBITRUM_SPEOLIA, types.ARBITRUM_ONE)

	L1GasOracleInL2 = map[types.Network]common.Address{
		types.OPTIMISM_MAINNET: common.HexToAddress("0x420000000000000000000000000000000000000F"),
	}
)

func GetNetWorkStack(network types.Network) types.NewWorkStack {
	if IsOpStackNetWork(network) {
		return types.OPSTACK
	}
	if IsArbNetWork(network) {
		return types.ARBSTACK
	}
	return types.DEFAULT_STACK
}

func IsTestNet(network types.Network) bool {
	return TestNetWork.Contains(network)
}
func IsOpStackNetWork(network types.Network) bool {
	return OpeStackNetWork.Contains(network)
}
func IsEthereumAdaptableNetWork(network types.Network) bool {
	return EthereumAdaptableNetWork.Contains(network)
}
func IsArbNetWork(network types.Network) bool {
	return ArbStackNetWork.Contains(network)
}
