package paymaster_data

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"math/big"
	"time"
)

type PaymasterDataInput struct {
	Paymaster                     common.Address
	ValidUntil                    *big.Int
	ValidAfter                    *big.Int
	ERC20Token                    common.Address
	ExchangeRate                  *big.Int
	PayType                       global_const.PayType
	EntryPointVersion             global_const.EntrypointVersion
	PaymasterVerificationGasLimit *big.Int
	PaymasterPostOpGasLimit       *big.Int
}

func NewPaymasterDataInput(strategy *model.Strategy) *PaymasterDataInput {
	now := time.Now()
	start := now.Add(-1 * time.Second)
	end := now.Add(5 * time.Minute)
	var tokenAddress string
	if strategy.GetPayType() == global_const.PayTypeERC20 {
		tokenAddress = config.GetTokenAddress(strategy.GetNewWork(), strategy.Erc20TokenType)

	} else {
		tokenAddress = global_const.DummyAddress.String()
		logrus.Debug("token address ", tokenAddress)
	}

	return &PaymasterDataInput{
		Paymaster:                     *strategy.GetPaymasterAddress(),
		ValidUntil:                    big.NewInt(end.Unix()),
		ValidAfter:                    big.NewInt(start.Unix()),
		ERC20Token:                    common.HexToAddress(tokenAddress),
		ExchangeRate:                  big.NewInt(0),
		PayType:                       strategy.GetPayType(),
		EntryPointVersion:             strategy.GetStrategyEntrypointVersion(),
		PaymasterVerificationGasLimit: global_const.DummyPaymasterOversimplificationBigint,
		PaymasterPostOpGasLimit:       global_const.DummyPaymasterPostoperativelyBigint,
	}
}
