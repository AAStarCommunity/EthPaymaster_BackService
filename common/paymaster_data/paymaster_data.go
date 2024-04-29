package paymaster_data

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type PaymasterData struct {
	Paymaster    common.Address
	ValidUntil   *big.Int
	ValidAfter   *big.Int
	ERC20Token   common.Address
	ExchangeRate *big.Int
}

func NewPaymasterDataInput(strategy *model.Strategy) *PaymasterData {
	start := strategy.ExecuteRestriction.EffectiveStartTime
	end := strategy.ExecuteRestriction.EffectiveEndTime
	tokenAddress := conf.GetTokenAddress(strategy.GetNewWork(), strategy.Erc20TokenType)
	return &PaymasterData{
		Paymaster:    *strategy.GetPaymasterAddress(),
		ValidUntil:   big.NewInt(end.Int64()),
		ValidAfter:   big.NewInt(start.Int64()),
		ERC20Token:   common.HexToAddress(tokenAddress),
		ExchangeRate: big.NewInt(0),
	}
}
