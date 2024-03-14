package paymaster_pay_type

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
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

func (e *Erc20PaymasterExecutor) GeneratePayMasterAndData(strategy *model.Strategy, userOp *model.UserOperation, gasResponse *model.ComputeGasResponse, extra map[string]any) (string, error) {
	validationGas := userOp.VerificationGasLimit.String()[0:16]
	postOPGas := userOp.CallGasLimit.String()[0:16]
	message := validationGas + postOPGas + string(strategy.PayType)

	signatureByte, err := utils.SignMessage("1d8a58126e87e53edc7b24d58d1328230641de8c4242c135492bf5560e0ff421", message)
	if err != nil {
		return "", err
	}
	signatureStr := hex.EncodeToString(signatureByte)
	message = message + signatureStr
	return message, nil
}
