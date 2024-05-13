package gas_executor

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"golang.org/x/xerrors"
	"math/big"
)

var (
	gasValidateFuncMap       = map[global_const.PayType]ValidatePaymasterGasFunc{}
	verifyingGasValidateFunc = func(userOp *user_op.UserOpInput, gasComputeResponse *model.ComputeGasResponse, strategy *model.Strategy) error {
		//Validate the account’s deposit in the entryPoint is high enough to cover the max possible cost (cover the already-done verification and max execution gas)
		// Paymaster check paymaster balance

		//check EntryPoint paymasterAddress balance
		balance, err := chain_service.GetPaymasterEntryPointBalance(strategy)
		if err != nil {
			return err
		}
		// if balance < 0
		if balance.Cmp(big.NewFloat(0)) < 0 {
			return xerrors.Errorf("paymaster EntryPoint balance < 0")
		}
		etherCost := gasComputeResponse.TotalGasDetail.MaxTxGasCostInEther
		if balance.Cmp(etherCost) < 0 {
			return xerrors.Errorf("paymaster EntryPoint Not Enough balance %s < %s", balance, etherCost)
		}
		return nil
	}
	erc20GasValidateFunc = func(userOp *user_op.UserOpInput, gasComputeResponse *model.ComputeGasResponse, strategy *model.Strategy) error {
		userOpValue := *userOp
		sender := userOpValue.Sender
		tokenBalance, getTokenBalanceErr := chain_service.GetAddressTokenBalance(strategy.GetNewWork(), *sender, strategy.Erc20TokenType)
		if getTokenBalanceErr != nil {
			return getTokenBalanceErr
		}
		tokenCost := gasComputeResponse.Erc20TokenCost
		bigFloatValue := new(big.Float).SetFloat64(tokenBalance)
		if bigFloatValue.Cmp(tokenCost) < 0 {
			return xerrors.Errorf("user Token Not Enough tokenBalance %s < tokenCost %s", tokenBalance, tokenCost)
		}
		return nil
	}
	superGasValidateFunc = func(userOp *user_op.UserOpInput, gasComputeResponse *model.ComputeGasResponse, strategy *model.Strategy) error {
		return xerrors.Errorf("never reach here")
	}
)

func init() {
	gasValidateFuncMap[global_const.PayTypeVerifying] = verifyingGasValidateFunc
	gasValidateFuncMap[global_const.PayTypeERC20] = erc20GasValidateFunc
	gasValidateFuncMap[global_const.PayTypeSuperVerifying] = superGasValidateFunc
}

func GetGasValidateFunc(payType global_const.PayType) ValidatePaymasterGasFunc {
	return gasValidateFuncMap[payType]
}

type ValidatePaymasterGasFunc = func(userOp *user_op.UserOpInput, gasComputeResponse *model.ComputeGasResponse, strategy *model.Strategy) error
