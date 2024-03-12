package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"golang.org/x/xerrors"
	"math/big"
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
	gasPriceInEtherStr := gasPriceInEther.Text('f', 18)
	result.MaxBasePriceGwei = gasPriceInGwei
	result.MaxBasePriceEther = &gasPriceInEtherStr

	priorityPriceInGwei := new(big.Float).SetInt(priorityPriceWei)
	priorityPriceInGwei.Quo(priorityPriceInGwei, GweiFactor)
	priorityPriceInEther := new(big.Float).SetInt(priorityPriceWei)
	priorityPriceInEther.Quo(priorityPriceInEther, EthWeiFactor)
	priorityPriceInEtherStr := priorityPriceInEther.Text('f', 18)
	result.MaxPriorityPriceGwei = priorityPriceInGwei
	result.MaxPriorityPriceEther = &priorityPriceInEtherStr
	result.MaxPriorityPriceGwei = priorityPriceInGwei
	result.MaxPriorityPriceEther = &priorityPriceInEtherStr
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
