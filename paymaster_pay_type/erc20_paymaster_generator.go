package paymaster_pay_type

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"golang.org/x/xerrors"
	"math/big"
)

type Erc20PaymasterExecutor struct {
}

func (e *Erc20PaymasterExecutor) ValidateGas(userOp *userop.BaseUserOp, gasResponse *model.ComputeGasResponse, strategy *model.Strategy) error {
	userOpValue := *userOp
	sender := userOpValue.GetSender()
	tokenBalance, getTokenBalanceErr := chain_service.GetAddressTokenBalance(strategy.GetNewWork(), *sender, strategy.GetUseToken())
	if getTokenBalanceErr != nil {
		return getTokenBalanceErr
	}
	tokenCost := gasResponse.TokenCost
	bigFloatValue := new(big.Float).SetFloat64(tokenBalance)
	if bigFloatValue.Cmp(tokenCost) < 0 {
		return xerrors.Errorf("user Token Not Enough tokenBalance %s < tokenCost %s", tokenBalance, tokenCost)
	}
	return nil
}
