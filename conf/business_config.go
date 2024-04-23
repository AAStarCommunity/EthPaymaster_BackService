package conf

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"os"
)

var basicConfig *BusinessConfig

func BusinessConfigInit(path string) {
	if path == "" {
		panic("pathParam is empty")
	}
	originConfig := initBusinessConfig(path)
	basicConfig = convertConfig(originConfig)
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
func initBusinessConfig(path string) *OriginBusinessConfig {
	var config OriginBusinessConfig
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("file not found: %s", path))
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

func GetSupportEntryPoints(network types.Network) (mapset.Set[string], error) {
	res, ok := basicConfig.SupportEntryPoint[network]
	if !ok {
		return nil, fmt.Errorf("network not found")
	}
	return res, nil
}
func GetSupportPaymaster(network types.Network) (mapset.Set[string], error) {
	res, ok := basicConfig.SupportPaymaster[network]
	if !ok {
		return nil, fmt.Errorf("network not found")
	}
	return res, nil
}

func GetTokenAddress(networkParam types.Network, tokenParam types.TokenType) string {
	networkConfig := basicConfig.NetworkConfigMap[networkParam]

	return networkConfig.TokenConfig[tokenParam]
}
func CheckEntryPointExist(network2 types.Network, address string) bool {
	entryPointSet := basicConfig.SupportEntryPoint[network2]
	entryPointSetValue := entryPointSet
	return entryPointSetValue.Contains(address)
}
func GetGasToken(networkParam types.Network) types.TokenType {
	networkConfig := basicConfig.NetworkConfigMap[networkParam]
	return networkConfig.GasToken
}

func GetChainId(newworkParam types.Network) string {
	networkConfig := basicConfig.NetworkConfigMap[newworkParam]
	return networkConfig.ChainId
}
func GetEthereumRpcUrl(network types.Network) string {
	networkConfig := basicConfig.NetworkConfigMap[network]
	return networkConfig.RpcUrl
}

var (
	TestNetWork = mapset.NewSet(
		types.EthereumSepolia, types.OptimismSepolia, types.ArbitrumSpeolia, types.ScrollSepolia, types.StarketSepolia, types.BaseSepolia)
	OpeStackNetWork = mapset.NewSet(
		types.OptimismMainnet, types.OptimismSepolia, types.BaseMainnet, types.BaseSepolia)
	EthereumAdaptableNetWork = mapset.NewSet(
		types.OptimismMainnet, types.OptimismSepolia, types.EthereumSepolia)
	ArbStackNetWork = mapset.NewSet(
		types.ArbitrumSpeolia, types.ArbitrumOne, types.ArbitrumNova)

	L1GasOracleInL2 = map[types.Network]common.Address{
		types.OptimismMainnet: common.HexToAddress("0x420000000000000000000000000000000000000F"),
		types.OptimismSepolia: common.HexToAddress("0x420000000000000000000000000000000000000F"),
		types.ScrollSepolia:   common.HexToAddress("0x5300000000000000000000000000000000000002"),
		types.ScrollMainnet:   common.HexToAddress("0x5300000000000000000000000000000000000002"),
	}
)

func GetNetWorkStack(network types.Network) types.NewWorkStack {
	if IsOpStackNetWork(network) {
		return types.OpStack
	}
	if IsArbNetWork(network) {
		return types.ArbStack
	}
	return types.DefaultStack
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
func GetSimulateEntryPointAddress(network types.Network) *common.Address {
	panic("implement me")
}
