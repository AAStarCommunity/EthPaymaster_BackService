package network

import (
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_common/contract/erc20"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/xerrors"
	"math/big"
)

var GweiFactor = new(big.Float).SetInt(big.NewInt(1e9))
var EthWeiFactor = new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

type EthereumExecutor struct {
	BaseExecutor
	Client  *ethclient.Client
	network types.Network
}

func GetEthereumExecutor(network types.Network) *EthereumExecutor {
	return nil
}

var TokenContractCache map[*common.Address]*contract_erc20.Contract
var ClientCache map[types.Network]*ethclient.Client

func init() {
	TokenContractCache = make(map[*common.Address]*contract_erc20.Contract)
}
func (executor EthereumExecutor) GetUserTokenBalance(userAddress common.Address, token types.TokenType) (*big.Int, error) {
	tokenAddress := conf.GetTokenAddress(executor.network, token)
	tokenInstance, err := executor.GetTokenContract(tokenAddress)
	if err != nil {
		return nil, err
	}
	return tokenInstance.BalanceOf(&bind.CallOpts{}, userAddress)
}
func (executor EthereumExecutor) CheckContractAddressAccess(contract *common.Address) (bool, error) {
	client := executor.Client

	code, err := client.CodeAt(context.Background(), *contract, nil)
	if err != nil {
		return false, err
	}
	if len(code) == 0 {
		return false, xerrors.Errorf("contract  [%s] address not exist in [%s] network", contract, executor.network)
	}
	return true, nil
}

func (executor EthereumExecutor) GetTokenContract(tokenAddress *common.Address) (*contract_erc20.Contract, error) {
	client := executor.Client
	contract, ok := TokenContractCache[tokenAddress]
	if !ok {
		erc20Contract, err := contract_erc20.NewContract(*tokenAddress, client)
		if err != nil {
			return nil, err
		}
		TokenContractCache[tokenAddress] = erc20Contract
		return erc20Contract, nil
	}
	return contract, nil
}

func (executor EthereumExecutor) EstimateUserOpGas(entrypointAddress *common.Address, userOpParam *userop.BaseUserOp) (uint64, error) {
	client := executor.Client
	userOpValue := *userOpParam
	userOpValue.GetSender()
	userOpValue.GetSender()
	res, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: *entrypointAddress,
		To:   userOpValue.GetSender(),
		Data: userOpValue.GetCallData(),
	})
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (executor EthereumExecutor) GetCurGasPrice() (*model.GasPrice, error) {

	client := executor.Client

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
