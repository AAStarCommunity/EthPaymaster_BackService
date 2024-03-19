package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"AAStarCommunity/EthPaymaster_BackService/service/chain_service"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"AAStarCommunity/EthPaymaster_BackService/service/gas_service"
	"AAStarCommunity/EthPaymaster_BackService/service/pay_service"
	"AAStarCommunity/EthPaymaster_BackService/service/validator_service"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/xerrors"
	"math/big"
	"strconv"
	"strings"
	"time"
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
	//The maxFeePerGas and maxPriorityFeePerGas are above a configurable minimum value that the client is willing to accept. At the minimum, they are sufficiently high to be included with the current block.basefee.

	//validate gas
	if err := gas_service.ValidateGas(userOp, gasResponse, strategy); err != nil {
		return nil, err
	}

	//pay
	payReceipt, payError := executePay(strategy, userOp, gasResponse)
	if payError != nil {
		return nil, payError
	}

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
		EntryPointAddress:  strategy.EntryPointAddress,
		PayMasterAddress:   strategy.PayMasterAddress,
		PayReceipt:         payReceipt,
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
func packUserOp(userOp *model.UserOperation) (string, []byte, error) {
	abiEncoder, err := abi.JSON(strings.NewReader(`[
    {
        "inputs": [
            {
                "components": [
                    {
                        "internalType": "address",
                        "name": "sender",
                        "type": "address"
                    },
                    {
                        "internalType": "uint256",
                        "name": "nonce",
                        "type": "uint256"
                    },
                    {
                        "internalType": "bytes",
                        "name": "initCode",
                        "type": "bytes"
                    },
                    {
                        "internalType": "bytes",
                        "name": "callData",
                        "type": "bytes"
                    },
                    {
                        "internalType": "uint256",
                        "name": "callGasLimit",
                        "type": "uint256"
                    },
                    {
                        "internalType": "uint256",
                        "name": "verificationGasLimit",
                        "type": "uint256"
                    },
                    {
                        "internalType": "uint256",
                        "name": "preVerificationGas",
                        "type": "uint256"
                    },
                    {
                        "internalType": "uint256",
                        "name": "maxFeePerGas",
                        "type": "uint256"
                    },
                    {
                        "internalType": "uint256",
                        "name": "maxPriorityFeePerGas",
                        "type": "uint256"
                    },
                    {
                        "internalType": "bytes",
                        "name": "paymasterAndData",
                        "type": "bytes"
                    },
                    {
                        "internalType": "bytes",
                        "name": "signature",
                        "type": "bytes"
                    }
                ],
                "internalType": "struct UserOperation",
                "name": "userOp",
                "type": "tuple"
            }
        ],
        "name": "UserOp",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
    }
	]`))
	if err != nil {
		return "", nil, err
	}
	method := abiEncoder.Methods["UserOp"]
	encoded, err := method.Inputs.Pack(userOp)

	if err != nil {
		return "", nil, err
	}
	//https://github.com/jayden-sudo/SoulWalletCore/blob/dc76bdb9a156d4f99ef41109c59ab99106c193ac/contracts/utils/CalldataPack.sol#L51-L65

	hexString := hex.EncodeToString(encoded)

	hexString = hexString[64:]
	hexLen := len(hexString)
	hexString = hexString[:hexLen-128]
	return hexString, encoded, nil
}

func UserOpHash(userOp *model.UserOperation, strategy *model.Strategy, validStart *big.Int, validEnd *big.Int) ([]byte, string, error) {
	packUserOpStr, _, err := packUserOp(userOp)
	if err != nil {
		return nil, "", err
	}
	//
	bytesTy, err := abi.NewType("bytes", "", nil)
	if err != nil {
		fmt.Println(err)
	}
	uint256Ty, err := abi.NewType("uint256", "", nil)
	if err != nil {
		fmt.Println(err)
	}

	addressTy, _ := abi.NewType("address", "", nil)
	arguments := abi.Arguments{
		{
			Type: bytesTy,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: addressTy,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: uint256Ty,
		},
	}
	chainId, err := chain_service.GetChainId(strategy.NetWork)
	if err != nil {
		return nil, "", err
	}
	packUserOpStrByteNew, _ := hex.DecodeString(packUserOpStr)
	chainId.Int64()
	bytesRes, err := arguments.Pack(packUserOpStrByteNew, chainId, common.HexToAddress(strategy.PayMasterAddress), userOp.Nonce, validStart, validEnd)
	if err != nil {
		return nil, "", err
	}
	encodeHash := crypto.Keccak256Hash(bytesRes)
	return encodeHash.Bytes(), hex.EncodeToString(bytesRes), nil

}

func getPayMasterAndData(strategy *model.Strategy, userOp *model.UserOperation, gasResponse *model.ComputeGasResponse) (string, string, error) {
	return generatePayMasterAndData(userOp, strategy)
}

func generatePayMasterAndData(userOp *model.UserOperation, strategy *model.Strategy) (string, string, error) {
	//v0.7 [0:20)paymaster address,[20:36)validation gas, [36:52)postop gas,[52:53)typeId,  [53:117)valid timestamp, [117:) signature
	//v0.6 [0:20)paymaster address,[20:22)payType, [22:86)start Time ,[86:150)typeId,  [53:117)valid timestamp, [117:) signature
	//validationGas := userOp.VerificationGasLimit.String()
	//postOPGas := userOp.CallGasLimit.String()
	validStart, validEnd := getValidTime()
	message := fmt.Sprintf("%s%s%s", strategy.PayMasterAddress, string(strategy.PayType), validStart+validEnd)
	signatureByte, err := SignPaymaster(userOp, strategy, validStart, validEnd)
	if err != nil {
		return "", "", err
	}
	signatureStr := hex.EncodeToString(signatureByte)
	message = message + signatureStr
	return message, signatureStr, nil
}

func SignPaymaster(userOp *model.UserOperation, strategy *model.Strategy, validStart string, validEnd string) ([]byte, error) {
	//string to int
	validStartInt, _ := strconv.ParseInt(validStart, 10, 64)
	validEndInt, _ := strconv.ParseInt(validEnd, 10, 64)
	userOpHash, _, err := UserOpHash(userOp, strategy, big.NewInt(validStartInt), big.NewInt(validEndInt))
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.HexToECDSA("1d8a58126e87e53edc7b24d58d1328230641de8c4242c135492bf5560e0ff421")
	if err != nil {
		return nil, err
	}
	signature, err := crypto.Sign(userOpHash, privateKey)
	return signature, err
}

func getValidTime() (string, string) {
	currentTime := time.Now()
	currentTimestamp := currentTime.Unix()
	futureTime := currentTime.Add(15 * time.Minute)
	futureTimestamp := futureTime.Unix()
	currentTimestampStr := strconv.FormatInt(currentTimestamp, 10)
	futureTimestampStr := strconv.FormatInt(futureTimestamp, 10)
	currentTimestampStrSupply := SupplyZero(currentTimestampStr, 64)
	futureTimestampStrSupply := SupplyZero(futureTimestampStr, 64)
	return currentTimestampStrSupply, futureTimestampStrSupply
}
func SupplyZero(prefix string, maxTo int) string {
	padding := maxTo - len(prefix)
	if padding > 0 {
		prefix = "0" + prefix
		prefix = fmt.Sprintf("%0*s", maxTo, prefix)
	}
	return prefix
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
