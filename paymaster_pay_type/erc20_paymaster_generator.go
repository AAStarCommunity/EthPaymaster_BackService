package paymaster_pay_type

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"golang.org/x/xerrors"
	"math/big"
)

type Erc20PaymasterExecutor struct {
}

func (e *Erc20PaymasterExecutor) ValidateGas(userOp *model.UserOperation, gasResponse *model.ComputeGasResponse, strategy *model.Strategy) error {
	tokenBalance, getTokenBalanceErr := chain_service.GetAddressTokenBalance(strategy.NetWork, userOp.Sender, strategy.Token)
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
