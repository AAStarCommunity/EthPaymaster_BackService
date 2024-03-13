package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"golang.org/x/xerrors"
	"math/big"
	"strings"
)

var GweiFactor = new(big.Float).SetInt(big.NewInt(1e9))
var EthWeiFactor = new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

func CheckContractAddressAccess(contract common.Address, chain types.Network) (bool, error) {
	if chain == "" {
		return false, xerrors.Errorf("chain can not be empty")
	}

	client, exist := EthCompatibleNetWorkClientMap[chain]
	if !exist {
		return false, xerrors.Errorf("chain Client [%s] not exist", chain)
	}
	code, err := client.CodeAt(context.Background(), contract, nil)
	if err != nil {
		return false, err
	}
	if len(code) == 0 {
		return false, xerrors.Errorf("contract  [%s] address not exist in [%s] network", contract, chain)
	}
	return true, nil
}

// GetGasPrice return gas price in wei, gwei, ether
func GetGasPrice(chain types.Network) (*model.GasPrice, error) {
	client, exist := EthCompatibleNetWorkClientMap[chain]
	if !exist {
		return nil, xerrors.Errorf("chain Client [%s] not exist", chain)
	}
	priceWei, priceWeiErr := client.SuggestGasPrice(context.Background())
	if priceWeiErr != nil {
		return nil, priceWeiErr
	}
	priorityPriceWei, tiperr := client.SuggestGasTipCap(context.Background())
	if tiperr != nil {
		return nil, tiperr
	}
	result := model.GasPrice{}
	result.MaxBasePriceWei = priceWei
	result.MaxPriorityPriceWei = priorityPriceWei

	gasPriceInGwei := new(big.Float).SetInt(priceWei)
	gasPriceInGwei.Quo(gasPriceInGwei, GweiFactor)
	gasPriceInEther := new(big.Float).SetInt(priceWei)
	gasPriceInEther.Quo(gasPriceInEther, EthWeiFactor)
	gasPriceInGweiFloat, _ := gasPriceInGwei.Float64()
	result.MaxBasePriceGwei = gasPriceInGweiFloat
	result.MaxBasePriceEther = gasPriceInEther

	priorityPriceInGwei := new(big.Float).SetInt(priorityPriceWei)
	priorityPriceInGwei.Quo(priorityPriceInGwei, GweiFactor)
	priorityPriceInEther := new(big.Float).SetInt(priorityPriceWei)
	priorityPriceInEther.Quo(priorityPriceInEther, EthWeiFactor)
	priorityPriceInGweiFloat, _ := priorityPriceInGwei.Float64()
	result.MaxPriorityPriceGwei = priorityPriceInGweiFloat
	result.MaxPriorityPriceEther = gasPriceInEther
	return &result, nil
}

func GetGas(netWork types.Network) (*big.Int, error) {
	client, exist := EthCompatibleNetWorkClientMap[netWork]
	if !exist {
		return nil, xerrors.Errorf("chain Client [%s] not exist", netWork)
	}
	head, erro := client.HeaderByNumber(context.Background(), nil)
	if erro != nil {
		return nil, erro
	}
	return head.BaseFee, nil
}
func GetPriorityFee(netWork types.Network) (*big.Int, *big.Float) {
	client, exist := EthCompatibleNetWorkClientMap[netWork]
	if !exist {
		return nil, nil
	}
	priceWei, _ := client.SuggestGasTipCap(context.Background())
	gasPriceInGwei := new(big.Float).SetInt(priceWei)
	gasPriceInGwei.Quo(gasPriceInGwei, GweiFactor)
	return priceWei, gasPriceInGwei
}

func GetEntryPointDeposit(entrypoint string, depositAddress string) uint256.Int {
	return uint256.Int{1}
}
func EstimateGasLimitAndCost(chain types.Network, msg ethereum.CallMsg) (uint64, error) {
	client, exist := EthCompatibleNetWorkClientMap[chain]
	if !exist {
		return 0, xerrors.Errorf("chain Client [%s] not exist", chain)
	}
	return client.EstimateGas(context.Background(), msg)
}
func GetAddressTokenBalance(network types.Network, address common.Address, token types.TokenType) ([]interface{}, error) {
	client, exist := EthCompatibleNetWorkClientMap[network]
	if !exist {
		return nil, xerrors.Errorf("chain Client [%s] not exist", network)
	}
	client.BalanceAt(context.Background(), address, nil)
	usdtContractAddress := common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7")
	//address := common.HexToAddress("0xDf7093eF81fa23415bb703A685c6331584D30177")
	const bananceABI = `[{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`
	usdtABI, jsonErr := abi.JSON(strings.NewReader(bananceABI))
	if jsonErr != nil {
		return nil, jsonErr
	}
	data, backErr := usdtABI.Pack("balanceOf", address)
	if backErr != nil {
		return nil, backErr

	}
	//usdtInstance, err := ethclient.NewContract(usdtContractAddress, usdtAbi, client)
	result, callErr := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &usdtContractAddress,
		Data: data,
	}, nil)
	if callErr != nil {
		return nil, callErr
	}
	var balanceResult, unpackErr = usdtABI.Unpack("balanceOf", result)
	if unpackErr != nil {
		return nil, unpackErr
	}
	//TODO get token balance
	return balanceResult, nil
}
