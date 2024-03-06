package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"AAStarCommunity/EthPaymaster_BackService/service/gas_service"
	"AAStarCommunity/EthPaymaster_BackService/service/pay_service"
	"AAStarCommunity/EthPaymaster_BackService/service/validator_service"
	"encoding/hex"
	"golang.org/x/xerrors"
)

func TryPayUserOpExecute(request *model.TryPayUserOpRequest) (*model.Result, error) {
	// validator
	if err := businessParamValidate(request); err != nil {
		return nil, err
	}
	userOp := request.UserOperation

	// getStrategy
	var strategy *model.Strategy
	if stg, err := strategyGenerate(request); err != nil {
		return nil, err
	} else if err = validator_service.ValidateStrategy(stg, &userOp); err != nil {
		return nil, err
	} else {
		strategy = stg
	}

	//base Strategy and UserOp computeGas
	gasResponse, gasComputeError := gas_service.ComputeGas(&userOp, strategy)
	if gasComputeError != nil {
		return nil, gasComputeError
	}

	//validate gas
	if err := gas_service.ValidateGas(&userOp, gasResponse); err != nil {
		return nil, err
	}

	//pay
	payReceipt, payError := executePay(strategy, &userOp, gasResponse)
	if payError != nil {
		return nil, payError
	}
	paymasterSignature := getPayMasterSignature(strategy, &userOp)
	var result = model.TryPayUserOpResponse{
		StrategyId:         strategy.Id,
		EntryPointAddress:  strategy.EntryPointAddress,
		PayMasterAddress:   strategy.PayMasterAddress,
		PayReceipt:         payReceipt,
		PayMasterSignature: paymasterSignature,
		GasInfo:            gasResponse,
	}

	return &model.Result{
		Code:    200,
		Data:    result,
		Message: "message",
		Cost:    "cost",
	}, nil
}

func businessParamValidate(request *model.TryPayUserOpRequest) error {
	if request.ForceStrategyId == "" && (request.ForceToken == "" || request.ForceNetwork == "") {
		return xerrors.Errorf("Token And Network Must Set When ForceStrategyId Is Empty")
	}
	//UserOp Validate
	if err := validator_service.ValidateUserOp(&request.UserOperation); err != nil {
		return err
	}
	if request.ForceEntryPointAddress != "" && request.ForceNetwork != "" {
		// check Address is available in NetWork
		if ok, err := chain_service.CheckContractAddressAccess(request.ForceEntryPointAddress, request.ForceNetwork); err != nil {
			return err
		} else if !ok {
			return xerrors.Errorf("ForceEntryPointAddress: [%s] not exist in [%s] network", request.ForceEntryPointAddress, request.ForceNetwork)
		}
	}
	return nil
}

func executePay(strategy *model.Strategy, userOp *model.UserOperationItem, gasResponse *model.ComputeGasResponse) (*model.PayReceipt, error) {
	//1.Recharge
	ethereumPayservice := pay_service.EthereumPayService{}
	if err := ethereumPayservice.Pay(); err != nil {
		return nil, err
	}
	//2.record account
	ethereumPayservice.RecordAccount()
	//3.return Receipt
	ethereumPayservice.GetReceipt()
	return &model.PayReceipt{
		TransactionHash: "0x110406d44ec1681fcdab1df2310181dee26ff43c37167b2c9c496b35cce69437",
		Sponsor:         "aastar",
	}, nil
}
func getPayMasterSignature(strategy *model.Strategy, userOp *model.UserOperationItem) string {
	signatureBytes, _ := utils.SignUserOp("1d8a58126e87e53edc7b24d58d1328230641de8c4242c135492bf5560e0ff421", userOp)
	return hex.EncodeToString(signatureBytes)
}

func strategyGenerate(request *model.TryPayUserOpRequest) (*model.Strategy, error) {
	if forceStrategyId := request.ForceStrategyId; forceStrategyId != "" {
		//force strategy
		if strategy := dashboard_service.GetStrategyById(forceStrategyId); strategy == nil {
			return nil, xerrors.Errorf("Not Support Strategy ID: [%w]", forceStrategyId)
		} else {
			return strategy, nil
		}
	}

	suitableStrategy, err := dashboard_service.GetSuitableStrategy(request.ForceEntryPointAddress, request.ForceNetwork, request.ForceToken) //TODO
	if err != nil {
		return nil, err
	}
	if suitableStrategy == nil {
		return nil, xerrors.Errorf("Empty Strategies")
	}
	return suitableStrategy, nil
}
