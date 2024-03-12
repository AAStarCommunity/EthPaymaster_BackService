package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"AAStarCommunity/EthPaymaster_BackService/paymaster_data_generator"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"AAStarCommunity/EthPaymaster_BackService/service/gas_service"
	"AAStarCommunity/EthPaymaster_BackService/service/pay_service"
	"AAStarCommunity/EthPaymaster_BackService/service/validator_service"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/xerrors"
)

func TryPayUserOpExecute(request *model.TryPayUserOpRequest) (*model.TryPayUserOpResponse, error) {
	// validator
	if err := businessParamValidate(request); err != nil {
		return nil, err
	}

	var strategy *model.Strategy
	// getStrategy
	strategy, generateErr := strategyGenerate(request)
	if generateErr != nil {
		return nil, generateErr
	}
	if strategy.EntryPointTag != types.EntrypointV06 {
		return nil, xerrors.Errorf("Not Support EntryPointTag: [%w]", strategy.EntryPointTag)
	}

	userOp, newUserOpError := model.NewUserOp(&request.UserOp)
	if newUserOpError != nil {
		return nil, newUserOpError
	}

	if err := validator_service.ValidateStrategy(strategy, userOp); err != nil {
		return nil, err
	}
	//recall simulate?
	//UserOp Validate
	//check nonce
	if err := validator_service.ValidateUserOp(userOp); err != nil {
		return nil, err
	}
	//base Strategy and UserOp computeGas
	gasResponse, gasComputeError := gas_service.ComputeGas(userOp, strategy)
	if gasComputeError != nil {
		return nil, gasComputeError
	}

	//validate gas
	if err := gas_service.ValidateGas(userOp, gasResponse, strategy); err != nil {
		return nil, err
	}

	//pay
	payReceipt, payError := executePay(strategy, userOp, gasResponse)
	if payError != nil {
		return nil, payError
	}
	paymasterSignature := getPayMasterSignature(strategy, userOp)

	var paymasterAndData []byte
	if paymasterAndDataRes, err := getPayMasterAndData(strategy, userOp, gasResponse, paymasterSignature); err != nil {
		return nil, err
	} else {
		paymasterAndData = paymasterAndDataRes
	}
	userOp.PaymasterAndData = paymasterAndData
	//validatePaymasterUserOp
	var result = &model.TryPayUserOpResponse{
		StrategyId:         strategy.Id,
		EntryPointAddress:  strategy.EntryPointAddress,
		PayMasterAddress:   strategy.PayMasterAddress,
		PayReceipt:         payReceipt,
		PayMasterSignature: paymasterSignature,
		PayMasterAndData:   utils.EncodeToStringWithPrefix(paymasterAndData),
		GasInfo:            gasResponse,
	}
	return result, nil
}

func businessParamValidate(request *model.TryPayUserOpRequest) error {
	if request.ForceStrategyId == "" && (request.ForceToken == "" || request.ForceNetwork == "") {
		return xerrors.Errorf("Token And Network Must Set When ForceStrategyId Is Empty")
	}
	if conf.Environment.IsDevelopment() && request.ForceNetwork != "" {
		if types.TestNetWork[request.ForceNetwork] {
			return xerrors.Errorf("Test Network Not Support")
		}
	}

	if request.ForceEntryPointAddress != "" && request.ForceNetwork != "" {
		// check Address is available in NetWork
		if ok, err := chain_service.CheckContractAddressAccess(common.HexToAddress(request.ForceEntryPointAddress), request.ForceNetwork); err != nil {
			return err
		} else if !ok {
			return xerrors.Errorf("ForceEntryPointAddress: [%s] not exist in [%s] network", request.ForceEntryPointAddress, request.ForceNetwork)
		}
	}
	return nil
}

func executePay(strategy *model.Strategy, userOp *model.UserOperation, gasResponse *model.ComputeGasResponse) (*model.PayReceipt, error) {
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
func getPayMasterSignature(strategy *model.Strategy, userOp *model.UserOperation) string {
	signatureBytes, _ := utils.SignUserOp("1d8a58126e87e53edc7b24d58d1328230641de8c4242c135492bf5560e0ff421", userOp)
	return hex.EncodeToString(signatureBytes)
}
func getPayMasterAndData(strategy *model.Strategy, userOp *model.UserOperation, gasResponse *model.ComputeGasResponse, paymasterSign string) ([]byte, error) {
	paymasterDataGenerator := paymaster_data_generator.GetPaymasterDataGenerator(strategy.PayType)
	if paymasterDataGenerator == nil {
		return nil, xerrors.Errorf("Not Support PayType: [%w]", strategy.PayType)
	}
	extra := make(map[string]any)
	extra["signature"] = paymasterSign
	return paymasterDataGenerator.GeneratePayMaster(strategy, userOp, gasResponse, extra)
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
