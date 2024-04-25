package gas_validate

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"golang.org/x/xerrors"
	"math/big"
)

var (
	GasValidateFuncMap = map[global_const.PayType]ValidatePaymasterGasFunc{}
)

func init() {
	GasValidateFuncMap[global_const.PayTypeVerifying] = VerifyingGasValidate()
	GasValidateFuncMap[global_const.PayTypeERC20] = Erc20GasValidate()
	GasValidateFuncMap[global_const.PayTypeSuperVerifying] = SuperGasValidate()
}

type ValidatePaymasterGasFunc = func(userOp *user_op.UserOpInput, gasComputeResponse *model.ComputeGasResponse, strategy *model.Strategy) error

func SuperGasValidate() ValidatePaymasterGasFunc {
	return func(userOp *user_op.UserOpInput, gasComputeResponse *model.ComputeGasResponse, strategy *model.Strategy) error {
		return xerrors.Errorf("never reach here")
	}
}
func Erc20GasValidate() ValidatePaymasterGasFunc {
	return func(userOp *user_op.UserOpInput, gasComputeResponse *model.ComputeGasResponse, strategy *model.Strategy) error {
		userOpValue := *userOp
		sender := userOpValue.Sender
		tokenBalance, getTokenBalanceErr := chain_service.GetAddressTokenBalance(strategy.GetNewWork(), *sender, strategy.GetUseToken())
		if getTokenBalanceErr != nil {
			return getTokenBalanceErr
		}
		tokenCost := gasComputeResponse.TokenCost
		bigFloatValue := new(big.Float).SetFloat64(tokenBalance)
		if bigFloatValue.Cmp(tokenCost) < 0 {
			return xerrors.Errorf("user Token Not Enough tokenBalance %s < tokenCost %s", tokenBalance, tokenCost)
		}
		return nil
	}
}
func VerifyingGasValidate() ValidatePaymasterGasFunc {
	return func(userOp *user_op.UserOpInput, gasComputeResponse *model.ComputeGasResponse, strategy *model.Strategy) error {
		//Validate the account’s deposit in the entryPoint is high enough to cover the max possible cost (cover the already-done verification and max execution gas)
		// Paymaster check paymaster balance

		//check EntryPoint paymasterAddress balance
		paymasterAddress := strategy.GetPaymasterAddress()
		tokenBalance, getTokenBalanceErr := chain_service.GetAddressTokenBalance(strategy.GetNewWork(), *paymasterAddress, strategy.GetUseToken())
		if getTokenBalanceErr != nil {
			return getTokenBalanceErr
		}
		tokenBalanceBigFloat := new(big.Float).SetFloat64(tokenBalance)
		if tokenBalanceBigFloat.Cmp(gasComputeResponse.TokenCost) > 0 {
			return xerrors.Errorf("paymaster Token Not Enough tokenBalance %s < tokenCost %s", tokenBalance, gasComputeResponse.TokenCost)
		}
		return nil
	}
}
