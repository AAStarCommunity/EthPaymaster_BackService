package config

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"os"
)

var basicConfig *BusinessConfig

func basicConfigInit(path string) {
	if path == "" {
		panic("pathParam is empty")
	}
	originConfig := initBusinessConfig(path)
	basicConfig = convertConfig(originConfig)
}

func convertConfig(originConfig *OriginBusinessConfig) *BusinessConfig {
	basic := &BusinessConfig{}
	basic.NetworkConfigMap = make(map[global_const.Network]NetWorkConfig)
	basic.SupportEntryPoint = make(map[global_const.Network]mapset.Set[string])
	basic.SupportPaymaster = make(map[global_const.Network]mapset.Set[string])
	for network, originNetWorkConfig := range originConfig.NetworkConfigMap {
		//TODO valid
		basic.NetworkConfigMap[network] = NetWorkConfig{
			ChainId:              originNetWorkConfig.ChainId,
			IsTest:               originNetWorkConfig.IsTest,
			TokenConfig:          originNetWorkConfig.TokenConfig,
			GasToken:             originNetWorkConfig.GasToken,
			V06EntryPointAddress: common.HexToAddress(originNetWorkConfig.V06EntryPointAddress),
			V07EntryPointAddress: common.HexToAddress(originNetWorkConfig.V07EntryPointAddress),
			V06PaymasterAddress:  common.HexToAddress(originNetWorkConfig.V06PaymasterAddress),
			V07PaymasterAddress:  common.HexToAddress(originNetWorkConfig.V07PaymasterAddress),
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
	NetworkConfigMap  map[global_const.Network]*OriginNetWorkConfig `json:"network_config"`
	SupportEntryPoint map[global_const.Network][]string             `json:"support_entrypoint"`
	SupportPaymaster  map[global_const.Network][]string             `json:"support_paymaster"`
}
type OriginNetWorkConfig struct {
	ChainId              string                            `json:"chain_id"`
	IsTest               bool                              `json:"is_test"`
	TokenConfig          map[global_const.TokenType]string `json:"token_config"`
	GasToken             global_const.TokenType            `json:"gas_token"`
	V06PaymasterAddress  string                            `json:"v06_paymaster_address"`
	V07PaymasterAddress  string                            `json:"v07_paymaster_address"`
	V06EntryPointAddress string                            `json:"v06_entrypoint_address"`
	V07EntryPointAddress string                            `json:"v07_entrypoint_address"`
	GasOracleAddress     string
}

type BusinessConfig struct {
	NetworkConfigMap  map[global_const.Network]NetWorkConfig      `json:"network_config"`
	SupportEntryPoint map[global_const.Network]mapset.Set[string] `json:"support_entrypoint"`
	SupportPaymaster  map[global_const.Network]mapset.Set[string] `json:"support_paymaster"`
}
type NetWorkConfig struct {
	ChainId              string                            `json:"chain_id"`
	IsTest               bool                              `json:"is_test"`
	TokenConfig          map[global_const.TokenType]string `json:"token_config"`
	GasToken             global_const.TokenType
	GasOracleAddress     common.Address
	V06PaymasterAddress  common.Address
	V07PaymasterAddress  common.Address
	V06EntryPointAddress common.Address
	V07EntryPointAddress common.Address
}

func GetSupportEntryPoints(network global_const.Network) (mapset.Set[string], error) {
	res, ok := basicConfig.SupportEntryPoint[network]
	if !ok {
		return nil, fmt.Errorf("network not found")
	}
	return res, nil
}
func GetSupportPaymaster(network global_const.Network) (mapset.Set[string], error) {
	res, ok := basicConfig.SupportPaymaster[network]
	if !ok {
		return nil, fmt.Errorf("network not found")
	}
	return res, nil
}

func GetTokenAddress(networkParam global_const.Network, tokenParam global_const.TokenType) string {
	networkConfig := basicConfig.NetworkConfigMap[networkParam]

	return networkConfig.TokenConfig[tokenParam]
}
func CheckEntryPointExist(network2 global_const.Network, address string) bool {
	entryPointSet := basicConfig.SupportEntryPoint[network2]
	entryPointSetValue := entryPointSet
	return entryPointSetValue.Contains(address)
}
func GetGasToken(networkParam global_const.Network) global_const.TokenType {
	networkConfig := basicConfig.NetworkConfigMap[networkParam]
	return networkConfig.GasToken
}

func GetChainId(networkParam global_const.Network) string {
	networkConfig := basicConfig.NetworkConfigMap[networkParam]
	return networkConfig.ChainId
}

func GetPaymasterAddress(network global_const.Network, version global_const.EntrypointVersion) common.Address {
	networkConfig := basicConfig.NetworkConfigMap[network]
	if version == global_const.EntrypointV07 {
		return networkConfig.V07PaymasterAddress
	}
	return networkConfig.V06PaymasterAddress
}

func GetEntrypointAddress(network global_const.Network, version global_const.EntrypointVersion) common.Address {
	networkConfig := basicConfig.NetworkConfigMap[network]
	if version == global_const.EntrypointV07 {
		return networkConfig.V07EntryPointAddress
	}
	return networkConfig.V06EntryPointAddress

}

var (
	TestNetWork = mapset.NewSet(
		global_const.EthereumSepolia, global_const.OptimismSepolia, global_const.ArbitrumSpeolia, global_const.ScrollSepolia, global_const.StarketSepolia, global_const.BaseSepolia)
	OpeStackNetWork = mapset.NewSet(
		global_const.OptimismMainnet, global_const.OptimismSepolia, global_const.BaseMainnet, global_const.BaseSepolia)
	EthereumAdaptableNetWork = mapset.NewSet(
		global_const.ArbitrumOne, global_const.ArbitrumNova, global_const.ArbitrumSpeolia,
		global_const.OptimismMainnet, global_const.OptimismSepolia, global_const.EthereumSepolia, global_const.EthereumMainnet, global_const.ScrollSepolia, global_const.ScrollMainnet, global_const.BaseMainnet, global_const.BaseSepolia)
	ArbStackNetWork = mapset.NewSet(
		global_const.ArbitrumSpeolia, global_const.ArbitrumOne, global_const.ArbitrumNova)

	L1GasOracleInL2 = map[global_const.Network]common.Address{
		global_const.OptimismMainnet: common.HexToAddress("0x420000000000000000000000000000000000000F"),
		global_const.OptimismSepolia: common.HexToAddress("0x420000000000000000000000000000000000000F"),
		global_const.BaseSepolia:     common.HexToAddress("0x420000000000000000000000000000000000000F"),
		global_const.BaseMainnet:     common.HexToAddress("0x420000000000000000000000000000000000000F"),
		global_const.ScrollSepolia:   common.HexToAddress("0x5300000000000000000000000000000000000002"),
		global_const.ScrollMainnet:   common.HexToAddress("0x5300000000000000000000000000000000000002"),
	}
	Disable1559Chain = mapset.NewSet(global_const.ScrollSepolia, global_const.ScrollMainnet)
)

func IsDisable1559Chain(network global_const.Network) bool {
	return Disable1559Chain.Contains(network)
}
func GetNetWorkStack(network global_const.Network) global_const.NewWorkStack {
	if IsOpStackNetWork(network) {
		return global_const.OpStack
	}
	if IsArbNetWork(network) {
		return global_const.ArbStack
	}
	return global_const.DefaultStack
}

func IsTestNet(network global_const.Network) bool {
	return TestNetWork.Contains(network)
}
func IsOpStackNetWork(network global_const.Network) bool {
	return OpeStackNetWork.Contains(network)
}
func IsEthereumAdaptableNetWork(network global_const.Network) bool {
	return EthereumAdaptableNetWork.Contains(network)
}
func IsArbNetWork(network global_const.Network) bool {
	return ArbStackNetWork.Contains(network)
}
