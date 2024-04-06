package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"AAStarCommunity/EthPaymaster_BackService/paymaster_pay_type"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"AAStarCommunity/EthPaymaster_BackService/service/gas_service"
	"AAStarCommunity/EthPaymaster_BackService/service/pay_service"
	"AAStarCommunity/EthPaymaster_BackService/service/validator_service"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/xerrors"
	"strings"
)

func TryPayUserOpExecute(request *model.TryPayUserOpRequest) (*model.TryPayUserOpResponse, error) {
	userOp, strategy, err := prepareExecute(request)
	if err != nil {
		return nil, err
	}
	gasResponse, err := estimateGas(userOp, strategy)
	if err != nil {
		return nil, err
	}
	payReceipt, err := executePay(strategy, userOp, gasResponse)
	if err != nil {
		return nil, err
	}
	result, err := postExecute(userOp, strategy, gasResponse)
	if err != nil {
		return nil, err
	}
	result.PayReceipt = payReceipt
	return result, nil
}

func prepareExecute(request *model.TryPayUserOpRequest) (*userop.BaseUserOp, *model.Strategy, error) {
	// validator
	if err := businessParamValidate(request); err != nil {
		return nil, nil, err
	}
	var strategy *model.Strategy
	// getStrategy
	strategy, generateErr := strategyGenerate(request)
	if generateErr != nil {
		return nil, nil, generateErr

	}

	userOp, err := userop.NewUserOp(&request.UserOp, strategy)
	if err != nil {
		return nil, nil, err

	}
	if err := validator_service.ValidateStrategy(strategy); err != nil {
		return nil, nil, err
	}
	if err := validator_service.ValidateUserOp(userOp); err != nil {
		return nil, nil, err
	}
	return userOp, strategy, nil
}

func estimateGas(userOp *userop.BaseUserOp, strategy *model.Strategy) (*model.ComputeGasResponse, error) {
	//base Strategy and UserOp computeGas
	gasResponse, gasComputeError := gas_service.ComputeGas(userOp, strategy)
	if gasComputeError != nil {
		return nil, gasComputeError
	}

	//The maxFeePerGas and maxPriorityFeePerGas are above a configurable minimum value that the client is willing to accept. At the minimum, they are sufficiently high to be included with the current block.basefee.

	//validate gas
	if err := gas_service.ValidateGas(userOp, gasResponse, strategy); err != nil {
		return nil, err
	}
	return gasResponse, nil
}

func postExecute(userOp *userop.BaseUserOp, strategy *model.Strategy, gasResponse *model.ComputeGasResponse) (*model.TryPayUserOpResponse, error) {
	var paymasterAndData string
	var paymasterSignature string
	if paymasterAndDataRes, paymasterSignatureRes, err := getPayMasterAndData(strategy, userOp, gasResponse); err != nil {
		return nil, err
	} else {
		paymasterAndData = paymasterAndDataRes
		paymasterSignature = paymasterSignatureRes
	}

	//validatePaymasterUserOp
	var result = &model.TryPayUserOpResponse{
		StrategyId:         strategy.Id,
		EntryPointAddress:  strategy.GetEntryPointAddress().String(),
		PayMasterAddress:   strategy.GetPaymasterAddress().String(),
		PayMasterSignature: paymasterSignature,
		PayMasterAndData:   paymasterAndData,
		GasInfo:            gasResponse,
	}
	return result, nil
}

func businessParamValidate(request *model.TryPayUserOpRequest) error {
	if request.ForceStrategyId == "" && (request.ForceToken == "" || request.ForceNetwork == "") {
		return xerrors.Errorf("Token And Network Must Set When ForceStrategyId Is Empty")
	}
	if conf.Environment.IsDevelopment() && request.ForceNetwork != "" {
		if network.TestNetWork[request.ForceNetwork] {
			return xerrors.Errorf(" %s not the Test Network ", request.ForceNetwork)
		}
	}
	exist := conf.CheckEntryPointExist(request.ForceNetwork, request.ForceEntryPointAddress)
	if !exist {
		return xerrors.Errorf("ForceEntryPointAddress: [%s] not exist in [%s] network", request.ForceEntryPointAddress, request.ForceNetwork)
	}
	return nil
}

func executePay(strategy *model.Strategy, userOp *userop.BaseUserOp, gasResponse *model.ComputeGasResponse) (*model.PayReceipt, error) {
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

func getPayMasterAndData(strategy *model.Strategy, userOp *userop.BaseUserOp, gasResponse *model.ComputeGasResponse) (string, string, error) {
	signatureByte, _, err := SignPaymaster(userOp, strategy)
	if err != nil {
		return "", "", err
	}
	signatureStr := hex.EncodeToString(signatureByte)
	dataGenerateFunc := paymaster_pay_type.GenerateFuncMap[strategy.GetPayType()]
	paymasterData, err := dataGenerateFunc(strategy, userOp, gasResponse)
	if err != nil {
		return "", "", err
	}
	paymasterDataResult := paymasterData + signatureStr
	return paymasterDataResult, signatureStr, err
}

func SignPaymaster(userOp *userop.BaseUserOp, strategy *model.Strategy) ([]byte, []byte, error) {
	userOpValue := *userOp
	userOpHash, _, err := userOpValue.GetUserOpHash(strategy)
	if err != nil {
		return nil, nil, err
	}
	hashToEthSignHash := utils.ToEthSignedMessageHash(userOpHash)

	privateKey, err := crypto.HexToECDSA("1d8a58126e87e53edc7b24d58d1328230641de8c4242c135492bf5560e0ff421")
	if err != nil {
		return nil, nil, err
	}

	signature, err := crypto.Sign(hashToEthSignHash, privateKey)

	signatureStr := hex.EncodeToString(signature)
	var signatureAfterProcess string

	if strings.HasSuffix(signatureStr, "00") {
		signatureAfterProcess = utils.ReplaceLastTwoChars(signatureStr, "1b")
	} else if strings.HasSuffix(signatureStr, "01") {
		signatureAfterProcess = utils.ReplaceLastTwoChars(signatureStr, "1c")
	} else {
		signatureAfterProcess = signatureStr
	}

	signatureAfterProcessByte, err := hex.DecodeString(signatureAfterProcess)
	if err != nil {
		return nil, nil, err
	}

	return signatureAfterProcessByte, userOpHash, err
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
