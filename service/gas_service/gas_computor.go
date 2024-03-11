package gas_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"encoding/hex"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/xerrors"
	"strconv"
)

func ComputeGas(userOp *model.UserOperationItem, strategy *model.Strategy) (*model.ComputeGasResponse, error) {
	priceInWei, gasPriceInGwei, gasPriceInEtherStr, getGasErr := chain_service.GetGasPrice(types.Sepolia)
	//TODO calculate the maximum possible fee the account needs to pay (based on validation and call gas limits, and current gas values)
	if getGasErr != nil {
		return nil, getGasErr
	}
	sender := common.HexToAddress(userOp.Sender)
	callData, _ := hex.DecodeString(userOp.CallData)
	//
	estimateCallGasLimit, _ := chain_service.EstimateGasLimitAndCost(strategy.NetWork, ethereum.CallMsg{
		From: common.HexToAddress(strategy.EntryPointAddress),
		To:   &sender,
		Data: callData,
	})
	userOpCallGasLimit, _ := strconv.ParseUint(userOp.CallGasLimit, 10, 64)
	if estimateCallGasLimit > userOpCallGasLimit {
		return nil, xerrors.Errorf("estimateCallGasLimit %d > userOpCallGasLimit %d", estimateCallGasLimit, userOpCallGasLimit)
	}

	// TODO get PaymasterCallGasLimit

	return &model.ComputeGasResponse{
		GasPriceInWei:   priceInWei.Uint64(),
		GasPriceInGwei:  gasPriceInGwei,
		GasPriceInEther: *gasPriceInEtherStr,
		CallGasLimit:    estimateCallGasLimit,
		TokenCost:       "0.0001",
		Network:         strategy.NetWork,
		Token:           strategy.Token,
		UsdCost:         "0.4",
		BlobEnable:      strategy.Enable4844,
	}, nil
}

func ValidateGas(userOp *model.UserOperationItem, gasComputeResponse *model.ComputeGasResponse) error {
	//1.if ERC20 check address balacnce
	//Validate the account’s deposit in the entryPoint is high enough to cover the max possible cost (cover the already-done verification and max execution gas)
	//2 if Paymaster check paymaster balance
	//The maxFeePerGas and maxPriorityFeePerGas are above a configurable minimum value that the client is willing to accept. At the minimum, they are sufficiently high to be included with the current block.basefee.
	return nil
}
