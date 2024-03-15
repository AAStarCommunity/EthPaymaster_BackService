package paymaster_pay_type

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"encoding/hex"
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

func (v VerifyingPaymasterExecutor) GeneratePayMasterAndData(strategy *model.Strategy, userOp *model.UserOperation, gasResponse *model.ComputeGasResponse, extra map[string]any) (string, error) {
	//[0:20)paymaster address,[20:36)validation gas, [36:52)postop gas,[52:53)typeId,  [53:117)valid timestamp, [117:) signature
	validationGas := userOp.VerificationGasLimit.String()
	postOPGas := userOp.CallGasLimit.String()
	message := strategy.PayMasterAddress + validationGas + postOPGas + string(strategy.PayType)
	//0000 timestamp /s (convert to hex)  64
	signatureByte, err := utils.SignMessage("1d8a58126e87e53edc7b24d58d1328230641de8c4242c135492bf5560e0ff421", message)
	if err != nil {
		return "", err
	}
	signatureStr := hex.EncodeToString(signatureByte)
	message = message + signatureStr
	return message, nil
}
