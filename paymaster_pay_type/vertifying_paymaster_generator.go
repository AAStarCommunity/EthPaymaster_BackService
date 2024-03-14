package paymaster_pay_type

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/xerrors"
	"math/big"
)

type VerifyingPaymasterExecutor struct {
}

func (v VerifyingPaymasterExecutor) ValidateGas(userOp *model.UserOperation, response *model.ComputeGasResponse, strategy *model.Strategy) error {
	//Validate the account’s deposit in the entryPoint is high enough to cover the max possible cost (cover the already-done verification and max execution gas)
	// Paymaster check paymaster balance

	//check EntryPoint paymasterAddress balance
	tokenBalance, getTokenBalanceErr := chain_service.GetAddressTokenBalance(strategy.NetWork, common.HexToAddress(strategy.PayMasterAddress), strategy.Token)
	if getTokenBalanceErr != nil {
		return getTokenBalanceErr
	}
	tokenBalanceBigFloat := new(big.Float).SetFloat64(tokenBalance)
	if tokenBalanceBigFloat.Cmp(response.TokenCost) > 0 {
		return xerrors.Errorf("paymaster Token Not Enough tokenBalance %s < tokenCost %s", tokenBalance, response.TokenCost)
	}
	return nil
}

func (v VerifyingPaymasterExecutor) GeneratePayMasterAndData(strategy *model.Strategy, userOp *model.UserOperation, gasResponse *model.ComputeGasResponse, extra map[string]any) ([]byte, error) {
	//verifying:[0-1]pay type，[1-21]paymaster address，[21-85]valid timestamp，[85-] signature
	signature, ok := extra["signature"]
	if !ok {
		return nil, xerrors.Errorf("signature not found")
	}
	res := "0x" + string(types.PayTypeVerifying) + strategy.PayMasterAddress + "" + signature.(string)
	return []byte(res), nil
}
