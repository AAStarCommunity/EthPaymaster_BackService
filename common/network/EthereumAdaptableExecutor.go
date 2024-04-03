package network

import (
	contract_erc20 "AAStarCommunity/EthPaymaster_BackService/common/contract/erc20"
	"AAStarCommunity/EthPaymaster_BackService/common/token"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var TokenContractCache map[*common.Address]*contract_erc20.Contract

func init() {
	TokenContractCache = make(map[*common.Address]*contract_erc20.Contract)
}
func GetUserTokenBalance(userAddress common.Address, token token.TokenType) {
}

func GetTokenContract(tokenAddress *common.Address, client *ethclient.Client) (*contract_erc20.Contract, error) {
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
