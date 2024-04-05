package network

import (
	contract_erc20 "AAStarCommunity/EthPaymaster_BackService/common/contract/erc20"
	"AAStarCommunity/EthPaymaster_BackService/common/tokens"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/xerrors"
	"math/big"
)

type EthereumExecutor struct {
	BaseExecutor
	Client  *ethclient.Client
	network Network
}

func GetEthereumExecutor(network Network) *EthereumExecutor {
	return nil
}

var TokenContractCache map[*common.Address]*contract_erc20.Contract
var ClientCache map[Network]*ethclient.Client

func init() {
	TokenContractCache = make(map[*common.Address]*contract_erc20.Contract)
}
func (executor EthereumExecutor) GetUserTokenBalance(userAddress common.Address, token tokens.TokenType) (*big.Int, error) {
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
