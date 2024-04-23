package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_common/contract/contract_paymaster_verifying_v07"
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_common/contract/paymater_verifying_erc20_v06"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"AAStarCommunity/EthPaymaster_BackService/paymaster_pay_type"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"AAStarCommunity/EthPaymaster_BackService/service/gas_service"
	"AAStarCommunity/EthPaymaster_BackService/service/pay_service"
	"AAStarCommunity/EthPaymaster_BackService/service/validator_service"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/xerrors"
	"math/big"
	"strings"
)

func TryPayUserOpExecute(request *model.UserOpRequest) (*model.TryPayUserOpResponse, error) {
	userOp, strategy, err := prepareExecute(request)
	if err != nil {
		return nil, err
	}
	gasResponse, paymasterUserOp, err := estimateGas(userOp, strategy)
	if err != nil {
		return nil, err
	}

	payReceipt, err := executePay(strategy, paymasterUserOp, gasResponse)
	if err != nil {
		return nil, err
	}
	result, err := postExecute(paymasterUserOp, strategy, gasResponse)
	if err != nil {
		return nil, err
	}
	result.PayReceipt = payReceipt
	return result, nil
}

//sub Function ---------

func prepareExecute(request *model.UserOpRequest) (*user_op.UserOpInput, *model.Strategy, error) {

	var strategy *model.Strategy

	strategy, generateErr := StrategyGenerate(request)
	if generateErr != nil {
		return nil, nil, generateErr
	}

	userOp, err := user_op.NewUserOp(&request.UserOp, strategy.GetStrategyEntryPointVersion())
	if err != nil {
		return nil, nil, err

	}
	if err := validator_service.ValidateStrategy(strategy); err != nil {
		return nil, nil, err
	}
	if err := validator_service.ValidateUserOp(userOp, strategy); err != nil {
		return nil, nil, err
	}
	return userOp, strategy, nil
}

func estimateGas(userOp *user_op.UserOpInput, strategy *model.Strategy) (*model.ComputeGasResponse, *user_op.UserOpInput, error) {
	//base Strategy and UserOp computeGas
	gasResponse, paymasterUserOp, gasComputeError := gas_service.ComputeGas(userOp, strategy)
	if gasComputeError != nil {
		return nil, nil, gasComputeError
	}
	//The maxFeePerGas and maxPriorityFeePerGas are above a configurable minimum value that the client is willing to accept. At the minimum, they are sufficiently high to be included with the current block.basefee.
	//validate gas
	if err := gas_service.ValidateGas(userOp, gasResponse, strategy); err != nil {
		return nil, nil, err
	}
	return gasResponse, paymasterUserOp, nil
}

func executePay(strategy *model.Strategy, userOp *user_op.UserOpInput, gasResponse *model.ComputeGasResponse) (*model.PayReceipt, error) {
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

func postExecute(userOp *user_op.UserOpInput, strategy *model.Strategy, gasResponse *model.ComputeGasResponse) (*model.TryPayUserOpResponse, error) {
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

func getPayMasterAndData(strategy *model.Strategy, userOp *user_op.UserOpInput, gasResponse *model.ComputeGasResponse) (string, string, error) {
	signatureByte, _, err := signPaymaster(userOp, strategy)
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

func signPaymaster(userOp *user_op.UserOpInput, strategy *model.Strategy) ([]byte, []byte, error) {
	userOpHash, _, err := GetUserOpHash(userOp, strategy)
	if err != nil {
		return nil, nil, err
	}
	signature, err := getUserOpHashSign(userOpHash)
	if err != nil {
		return nil, nil, err
	}

	return signature, userOpHash, err
}

func getUserOpHashSign(userOpHash []byte) ([]byte, error) {
	privateKey, err := crypto.HexToECDSA("1d8a58126e87e53edc7b24d58d1328230641de8c4242c135492bf5560e0ff421")
	if err != nil {
		return nil, err
	}

	signature, err := crypto.Sign(userOpHash, privateKey)
	signatureStr := hex.EncodeToString(signature)
	var signatureAfterProcess string
	if strings.HasSuffix(signatureStr, "00") {
		signatureAfterProcess = utils.ReplaceLastTwoChars(signatureStr, "1b")
	} else if strings.HasSuffix(signatureStr, "01") {
		signatureAfterProcess = utils.ReplaceLastTwoChars(signatureStr, "1c")
	} else {
		signatureAfterProcess = signatureStr
	}
	return hex.DecodeString(signatureAfterProcess)
}

func StrategyGenerate(request *model.UserOpRequest) (*model.Strategy, error) {
	if forceStrategyId := request.ForceStrategyId; forceStrategyId != "" {
		//force strategy
		if strategy := dashboard_service.GetStrategyById(forceStrategyId); strategy == nil {
			return nil, xerrors.Errorf("Not Support Strategy ID: [%w]", forceStrategyId)
		} else {
			return strategy, nil
		}
	}

	suitableStrategy, err := dashboard_service.GetSuitableStrategy(request.ForceEntryPointAddress, request.ForceNetwork, types.PayTypeSuperVerifying) //TODO
	if err != nil {
		return nil, err
	}
	if suitableStrategy == nil {
		return nil, xerrors.Errorf("Empty Strategies")
	}
	return suitableStrategy, nil
}

func GetUserOpHash(userOp *user_op.UserOpInput, strategy *model.Strategy) ([]byte, string, error) {
	version := strategy.GetStrategyEntryPointVersion()
	executor := network.GetEthereumExecutor(strategy.GetNewWork())
	erc20Token := common.HexToAddress("0x")
	paytype := strategy.GetPayType()
	if paytype == types.PayTypeERC20 {
		tokenType := strategy.GetUseToken()
		tokenAddress := conf.GetTokenAddress(strategy.GetNewWork(), tokenType)
		erc20Token = common.HexToAddress(tokenAddress)
	}

	if version == types.EntrypointV06 {
		contract, err := executor.GetPaymasterErc20AndVerifyV06(strategy.GetPaymasterAddress())
		if err != nil {
			return nil, "", err
		}
		hash, err := contract.GetHash(&bind.CallOpts{}, paymater_verifying_erc20_v06.UserOperation{
			Sender:               *userOp.Sender,
			Nonce:                userOp.Nonce,
			InitCode:             userOp.InitCode,
			CallData:             userOp.CallData,
			CallGasLimit:         userOp.CallGasLimit,
			VerificationGasLimit: userOp.VerificationGasLimit,
			PreVerificationGas:   userOp.PreVerificationGas,
			MaxFeePerGas:         userOp.MaxFeePerGas,
			MaxPriorityFeePerGas: userOp.MaxPriorityFeePerGas,
			PaymasterAndData:     userOp.PaymasterAndData,
			Signature:            userOp.Signature,
		}, strategy.ExecuteRestriction.EffectiveEndTime, strategy.ExecuteRestriction.EffectiveStartTime, erc20Token, big.NewInt(0))
		if err != nil {
			return nil, "", err
		}
		return hash[:], "", nil
	} else if version == types.EntryPointV07 {
		if paytype == types.PayTypeVerifying {
			contract, err := executor.GetPaymasterVerifyV07(strategy.GetPaymasterAddress())
			if err != nil {
				return nil, "", err
			}
			hash, err := contract.GetHash(&bind.CallOpts{}, contract_paymaster_verifying_v07.PackedUserOperation{
				Sender:   *userOp.Sender,
				Nonce:    userOp.Nonce,
				InitCode: userOp.InitCode,
				CallData: userOp.CallData,
				//TODO
			}, strategy.ExecuteRestriction.EffectiveEndTime, strategy.ExecuteRestriction.EffectiveStartTime)
			if err != nil {
				return nil, "", err
			}
			return hash[:], "", nil
		} else if paytype == types.PayTypeERC20 {
			//TODO
			panic("implement me")
			//contract, err := executor.GetPaymasterErc20V07(strategy.GetPaymasterAddress())
			//if err != nil {
			//	return nil, "", err
			//}
			//hash, err := contract.GetHash(&bind.CallOpts{}, contract_paymaster_e_v07.PackedUserOperation{}, strategy.ExecuteRestriction.EffectiveEndTime, strategy.ExecuteRestriction.EffectiveStartTime, erc20Token, big.NewInt(0))
			//if err != nil {
			//	return nil, "", err
			//}
			//return hash[:], "", nil

		} else {
			return nil, "", xerrors.Errorf("paytype %s not support", paytype)
		}
	} else {
		return nil, "", xerrors.Errorf("entrypoint version %s not support", version)
	}
	//paymasterGasValue := userOp.PaymasterPostOpGasLimit.Text(20) + userOp.PaymasterVerificationGasLimit.Text(20)
	//byteRes, err := UserOpV07GetHashArguments.Pack(userOp.Sender, userOp.Nonce, crypto.Keccak256(userOp.InitCode),
	//	crypto.Keccak256(userOp.CallData), userOp.AccountGasLimit,
	//	paymasterGasValue, userOp.PreVerificationGas, userOp.GasFees, conf.GetChainId(strategy.GetNewWork()), strategy.GetPaymasterAddress())
	//if err != nil {
	//	return nil, "", err
	//}
	//userOpHash := crypto.Keccak256(byteRes)
	//afterProcessUserOphash := utils.ToEthSignedMessageHash(userOpHash)
	//return afterProcessUserOphash, hex.EncodeToString(byteRes), nil

	// V06
	//packUserOpStr, _, err := packUserOpV6ForUserOpHash(userOp)
	//if err != nil {
	//	return nil, "", err
	//}
	//
	//packUserOpStrByteNew, err := hex.DecodeString(packUserOpStr)
	//if err != nil {
	//	return nil, "", err
	//}
	//
	//bytesRes, err := userOPV06GetHashArguments.Pack(packUserOpStrByteNew, conf.GetChainId(strategy.GetNewWork()), strategy.GetPaymasterAddress(), userOp.Nonce, strategy.ExecuteRestriction.EffectiveStartTime, strategy.ExecuteRestriction.EffectiveEndTime)
	//if err != nil {
	//	return nil, "", err
	//}
	//
	//userOpHash := crypto.Keccak256(bytesRes)
	//afterProcessUserOphash := utils.ToEthSignedMessageHash(userOpHash)
	//return afterProcessUserOphash, hex.EncodeToString(bytesRes), nil
	//TODO
	panic("implement me")

}
