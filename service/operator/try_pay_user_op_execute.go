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
                        "name": "Sender",
                        "type": "address"
                    },
                    {
                        "internalType": "uint256",
                        "name": "Nonce",
                        "type": "uint256"
                    },
                    {
                        "internalType": "bytes",
                        "name": "InitCode",
                        "type": "bytes"
                    },
                    {
                        "internalType": "bytes",
                        "name": "CallData",
                        "type": "bytes"
                    },
                    {
                        "internalType": "uint256",
                        "name": "CallGasLimit",
                        "type": "uint256"
                    },
                    {
                        "internalType": "uint256",
                        "name": "VerificationGasLimit",
                        "type": "uint256"
                    },
                    {
                        "internalType": "uint256",
                        "name": "PreVerificationGas",
                        "type": "uint256"
                    },
                    {
                        "internalType": "uint256",
                        "name": "MaxFeePerGas",
                        "type": "uint256"
                    },
                    {
                        "internalType": "uint256",
                        "name": "MaxPriorityFeePerGas",
                        "type": "uint256"
                    },
                    {
                        "internalType": "bytes",
                        "name": "PaymasterAndData",
                        "type": "bytes"
                    },
                    {
                        "internalType": "bytes",
                        "name": "Signature",
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
	//TODO disgusting logic
	userOp.PaymasterAndData = []byte("0xE99c4Db5E360B8c84bF3660393CB2A85c3029b4400000000000000000000000000000000000000000000000000000000171004449600000000000000000000000000000000000000000000000000000017415804969e46721fc1938ac427add8a9e0d5cba2be5b17ccda9b300d0d3eeaff1904dfc23e276abd1ba6e3e269ec6aa36fe6a2442c18d167b53d7f9f0d1b3ebe80b09a6200")
	encoded, err := method.Inputs.Pack(userOp)

	if err != nil {
		return "", nil, err
	}
	//https://github.com/jayden-sudo/SoulWalletCore/blob/dc76bdb9a156d4f99ef41109c59ab99106c193ac/contracts/utils/CalldataPack.sol#L51-L65
	hexString := hex.EncodeToString(encoded)

	//1. 从 63*10+ 1 ～64*10获取
	hexString = hexString[64:]
	//hexLen := len(hexString)
	subIndex := GetIndex(hexString)
	hexString = hexString[:subIndex]
	//fmt.Printf("subIndex: %d\n", subIndex)
	return hexString, encoded, nil
}
func GetIndex(hexString string) int64 {
	//1. 从 63*10+ 1 ～64*10获取

	indexPre := hexString[576:640]
	indePreInt, _ := strconv.ParseInt(indexPre, 16, 64)
	result := indePreInt * 2
	return result
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
	uint48Ty, err := abi.NewType("uint48", "", nil)

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
			Type: uint48Ty,
		},
		{
			Type: uint48Ty,
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

// 1710044496
// 1741580496
func getValidTime() (string, string) {
	//currentTime := time.Nsow()
	//currentTimestamp := 1710044496
	//futureTime := currentTime.Add(15 * time.Minute)
	//futureTimestamp := futureTime.Unix()
	currentTimestampStr := strconv.FormatInt(1710044496, 10)
	futureTimestampStr := strconv.FormatInt(1741580496, 10)
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
