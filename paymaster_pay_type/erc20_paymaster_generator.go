package paymaster_pay_type

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"encoding/hex"
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

func (e *Erc20PaymasterExecutor) GeneratePayMasterAndData(strategy *model.Strategy, userOp *model.UserOperation, gasResponse *model.ComputeGasResponse, extra map[string]any) ([]byte, error) {
	//ERC20:[0-1]pay type，[1-21]paymaster address，[21-53]token Amount
	//tokenCost := gasResponse.TokenCost.Float64()
	res := "0x" + string(types.PayTypeERC20) + strategy.PayMasterAddress
	//TODO implement me
	return hex.DecodeString(res)
}
