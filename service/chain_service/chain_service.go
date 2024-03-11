package chain_service

import (
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

func CheckContractAddressAccess(contract string, chain types.Network) (bool, error) {
	if chain == "" {
		return false, xerrors.Errorf("chain can not be empty")
	}
	contractAddress := common.HexToAddress(contract)

	client, exist := NetWorkClientMap[chain]
	if !exist {
		return false, xerrors.Errorf("chain Client [%s] not exist", chain)
	}
	code, err := client.CodeAt(context.Background(), contractAddress, nil)
	if err != nil {
		return false, err
	}
	if len(code) == 0 {
		return false, xerrors.Errorf("contract  [%s] address not exist in [%s] network", contract, chain)
	}
	return true, nil
}

// GetGasPrice return gas price in wei, gwei, ether
func GetGasPrice(chain types.Network) (*big.Int, *big.Float, *string, error) {
	client, exist := NetWorkClientMap[chain]
	if !exist {
		return nil, nil, nil, xerrors.Errorf("chain Client [%s] not exist", chain)
	}
	priceWei, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}

	gasPriceInGwei := new(big.Float).SetInt(priceWei)
	gasPriceInGwei.Quo(gasPriceInGwei, GweiFactor)

	gasPriceInEther := new(big.Float).SetInt(priceWei)
	gasPriceInEther.Quo(gasPriceInEther, EthWeiFactor)
	gasPriceInEtherStr := gasPriceInEther.Text('f', 18)
	return priceWei, gasPriceInGwei, &gasPriceInEtherStr, nil
}

func GetEntryPointDeposit(entrypoint string, depositAddress string) uint256.Int {
	return uint256.Int{1}
}
func EstimateGasLimitAndCost(chain types.Network, msg ethereum.CallMsg) (uint64, error) {
	client, exist := NetWorkClientMap[chain]
	if !exist {
		return 0, xerrors.Errorf("chain Client [%s] not exist", chain)
	}
	return client.EstimateGas(context.Background(), msg)
}
