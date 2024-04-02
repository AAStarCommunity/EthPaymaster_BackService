package chain_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/erc20_token"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"golang.org/x/xerrors"
	"math"
	"math/big"
	"strings"
)

var GweiFactor = new(big.Float).SetInt(big.NewInt(1e9))
var EthWeiFactor = new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

const balanceOfAbi = `[{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`

var TokenAddressMap map[network.Network]*map[erc20_token.TokenType]common.Address

func init() {
	TokenAddressMap = map[network.Network]*map[erc20_token.TokenType]common.Address{
		network.Ethereum: {
			erc20_token.ETH: common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7"),
		},
		network.Sepolia: {
			erc20_token.USDT: common.HexToAddress("0xaa8e23fb1079ea71e0a56f48a2aa51851d8433d0"),
			erc20_token.USDC: common.HexToAddress("0x1c7d4b196cb0c7b01d743fbc6116a902379c7238"),
		},
	}
}
func CheckContractAddressAccess(contract *common.Address, chain network.Network) (bool, error) {
	if chain == "" {
		return false, xerrors.Errorf("chain can not be empty")
	}
	client, exist := EthCompatibleNetWorkClientMap[chain]
	if !exist {
		return false, xerrors.Errorf("chain Client [%s] not exist", chain)
	}
	code, err := client.CodeAt(context.Background(), *contract, nil)
	if err != nil {
		return false, err
	}
	if len(code) == 0 {
		return false, xerrors.Errorf("contract  [%s] address not exist in [%s] network", contract, chain)
	}
	return true, nil
}

// GetGasPrice return gas price in wei, gwei, ether
func GetGasPrice(chain network.Network) (*model.GasPrice, error) {
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

func GetGas(netWork network.Network) (*big.Int, error) {
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
func GetPriorityFee(netWork network.Network) (*big.Int, *big.Float) {
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
func EstimateGasLimitAndCost(chain network.Network, msg ethereum.CallMsg) (uint64, error) {
	client, exist := EthCompatibleNetWorkClientMap[chain]
	if !exist {
		return 0, xerrors.Errorf("chain Client [%s] not exist", chain)
	}
	return client.EstimateGas(context.Background(), msg)
}
func GetAddressTokenBalance(network network.Network, address common.Address, token erc20_token.TokenType) (float64, error) {
	client, exist := EthCompatibleNetWorkClientMap[network]
	if !exist {
		return 0, xerrors.Errorf("chain Client [%s] not exist", network)
	}
	if token == erc20_token.ETH {
		res, err := client.BalanceAt(context.Background(), address, nil)
		if err != nil {
			return 0, err
		}
		bananceV := float64(res.Int64()) * math.Pow(10, -18)
		return bananceV, nil
	}

	tokenContractAddress := (*TokenAddressMap[network])[token]
	usdtABI, jsonErr := abi.JSON(strings.NewReader(balanceOfAbi))
	if jsonErr != nil {
		return 0, jsonErr
	}
	data, backErr := usdtABI.Pack("balanceOf", address)
	if backErr != nil {
		return 0, backErr
	}
	result, callErr := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &tokenContractAddress,
		Data: data,
	}, nil)
	if callErr != nil {
		return 0, callErr
	}

	var balanceResult *big.Int
	unpackErr := usdtABI.UnpackIntoInterface(&balanceResult, "balanceOf", result)
	if unpackErr != nil {
		return 0, unpackErr
	}
	balanceResultFloat := float64(balanceResult.Int64()) * math.Pow(10, -6)

	return balanceResultFloat, nil

}
func GetChainId(chain network.Network) (*big.Int, error) {
	client, exist := EthCompatibleNetWorkClientMap[chain]
	if !exist {
		return nil, xerrors.Errorf("chain Client [%s] not exist", chain)
	}
	return client.ChainID(context.Background())
}
