package conf

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"os"
)

var basicConfig *BusinessConfig
var signerConfig = make(SignerConfigMap)

type SignerConfigMap map[global_const.Network]*global_const.EOA

func BusinessConfigInit(path string) {
	if path == "" {
		panic("pathParam is empty")
	}
	originConfig := initBusinessConfig(path)
	basicConfig = convertConfig(originConfig)

}
func GetSigner(network global_const.Network) *global_const.EOA {
	return signerConfig[network]
}

func convertConfig(originConfig *OriginBusinessConfig) *BusinessConfig {
	basic := &BusinessConfig{}
	basic.NetworkConfigMap = make(map[global_const.Network]NetWorkConfig)
	basic.SupportEntryPoint = make(map[global_const.Network]mapset.Set[string])
	basic.SupportPaymaster = make(map[global_const.Network]mapset.Set[string])
	for network, originNetWorkConfig := range originConfig.NetworkConfigMap {
		//TODO valid
		basic.NetworkConfigMap[network] = NetWorkConfig{
			ChainId:     originNetWorkConfig.ChainId,
			IsTest:      originNetWorkConfig.IsTest,
			RpcUrl:      fmt.Sprintf("%s/%s", originNetWorkConfig.RpcUrl, originNetWorkConfig.ApiKey),
			TokenConfig: originNetWorkConfig.TokenConfig,
			GasToken:    originNetWorkConfig.GasToken,
		}
		paymasterArr := originConfig.SupportPaymaster[network]
		paymasterSet := mapset.NewSet[string]()
		paymasterSet.Append(paymasterArr...)
		basic.SupportPaymaster[network] = paymasterSet

		entryPointArr := originConfig.SupportEntryPoint[network]
		entryPointSet := mapset.NewSet[string]()
		entryPointSet.Append(entryPointArr...)
		basic.SupportEntryPoint[network] = entryPointSet
		//TODO starknet
		eoa, err := global_const.NewEoa(originNetWorkConfig.SignerKey)
		if err != nil {
			panic(fmt.Sprintf("signer key error: %s", err))
		}
		signerConfig[network] = eoa
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
	ChainId          string                            `json:"chain_id"`
	IsTest           bool                              `json:"is_test"`
	RpcUrl           string                            `json:"rpc_url"`
	ApiKey           string                            `json:"api_key"`
	SignerKey        string                            `json:"signer_key"`
	TokenConfig      map[global_const.TokenType]string `json:"token_config"`
	GasToken         global_const.TokenType            `json:"gas_token"`
	GasOracleAddress string
}

type BusinessConfig struct {
	NetworkConfigMap  map[global_const.Network]NetWorkConfig      `json:"network_config"`
	SupportEntryPoint map[global_const.Network]mapset.Set[string] `json:"support_entrypoint"`
	SupportPaymaster  map[global_const.Network]mapset.Set[string] `json:"support_paymaster"`
}
type NetWorkConfig struct {
	ChainId          string                            `json:"chain_id"`
	IsTest           bool                              `json:"is_test"`
	RpcUrl           string                            `json:"rpc_url"`
	TokenConfig      map[global_const.TokenType]string `json:"token_config"`
	GasToken         global_const.TokenType
	GasOracleAddress common.Address
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

func GetChainId(newworkParam global_const.Network) string {
	networkConfig := basicConfig.NetworkConfigMap[newworkParam]
	return networkConfig.ChainId
}
func GetEthereumRpcUrl(network global_const.Network) string {
	networkConfig := basicConfig.NetworkConfigMap[network]
	return networkConfig.RpcUrl
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
func GetSimulateEntryPointAddress(network global_const.Network) *common.Address {
	panic("implement me")
}
