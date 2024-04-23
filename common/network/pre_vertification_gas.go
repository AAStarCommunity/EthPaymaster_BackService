package network

import (
	"AAStarCommunity/EthPaymaster_BackService/common/arbitrum"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"golang.org/x/xerrors"
	"math"
	"math/big"
)

var preVerificationGasFuncMap = map[types.NewWorkStack]PreVerificationGasFunc{}

type PreVerificationGasFunc = func(op *userop.UserOpInput, strategy *model.Strategy, gasFeeResult *model.GasPrice) (*big.Int, error)

func init() {
	preVerificationGasFuncMap[types.ArbStack] = ArbitrumPreVerificationGasFunc()
	preVerificationGasFuncMap[types.DefaultStack] = DefaultPreVerificationGasFunc()
	preVerificationGasFuncMap[types.OpStack] = OPStackPreVerificationGasFunc()
}
func GetPreVerificationGasFunc(stack types.NewWorkStack) (PreVerificationGasFunc, error) {
	function, ok := preVerificationGasFuncMap[stack]
	if !ok {
		return nil, xerrors.Errorf("stack %s not support", stack)
	}
	return function, nil
}

// https://medium.com/offchainlabs/understanding-arbitrum-2-dimensional-fees-fd1d582596c9.
// https://docs.arbitrum.io/build-decentralized-apps/nodeinterface/reference
func ArbitrumPreVerificationGasFunc() PreVerificationGasFunc {
	return func(op *userop.UserOpInput, strategy *model.Strategy, gasFeeResult *model.GasPrice) (*big.Int, error) {
		base, err := getBasicPreVerificationGas(op, strategy)
		if err != nil {
			return nil, err
		}
		executor := GetEthereumExecutor(strategy.GetNewWork())
		estimateOutPut, err := arbitrum.GetEstimateL1ComponentMethod(executor.Client)
		if err != nil {
			return nil, err
		}
		big.NewInt(0).Add(base, big.NewInt(int64(estimateOutPut.GasEstimateForL1)))
		return big.NewInt(0), nil
	}
}
func DefaultPreVerificationGasFunc() PreVerificationGasFunc {
	return func(op *userop.UserOpInput, strategy *model.Strategy, gasFeeResult *model.GasPrice) (*big.Int, error) {
		return getBasicPreVerificationGas(op, strategy)
	}
}

// OPStackPreVerificationGasFunc
// The L1 data fee is paid based on the current Ethereum gas price as tracked within the GasPriceOracle smart contract. This gas price is updated automatically by the OP Mainnet protocol.
// https://docs.optimism.io/builders/app-developers/transactions/estimates#execution-gas-fee
// https://docs.optimism.io/stack/transactions/fees#the-l1-data-fee
func OPStackPreVerificationGasFunc() PreVerificationGasFunc {
	return func(op *userop.UserOpInput, strategy *model.Strategy, gasFeeResult *model.GasPrice) (*big.Int, error) {
		basicGas, err := getBasicPreVerificationGas(op, strategy)
		if err != nil {
			return nil, err
		}
		executor := GetEthereumExecutor(strategy.GetNewWork())
		data, err := getInputData(op)
		if err != nil {
			return nil, err
		}
		l1DataFee, err := executor.GetL1DataFee(data)
		if err != nil {
			return nil, err
		}
		l2MaxFee := gasFeeResult.MaxFeePerGas
		l2priorityFee := big.NewInt(0).Add(gasFeeResult.MaxFeePerGas, gasFeeResult.BaseFee)
		// use smaller one
		var l2Price *big.Int
		if utils.LeftIsLessTanRight(l2MaxFee, l2priorityFee) {
			l2Price = l2MaxFee
		} else {
			l2Price = l2priorityFee
		}
		//Return static + L1 buffer as PVG. L1 buffer is equal to L1Fee/L2Price.
		return big.NewInt(0).Add(basicGas, big.NewInt(0).Div(l1DataFee, l2Price)), nil
	}
}

/**
 * calculate the preVerificationGas of the given UserOperation
 * preVerificationGas (by definition) is the cost overhead that can't be calculated on-chain.
 * it is based on parameters that are defined by the Ethereum protocol for external transactions.
 * @param userOp filled userOp to calculate. The only possible missing fields can be the signature and preVerificationGas itself
 * @param overheads gas overheads to use, to override the default values
 */
func getBasicPreVerificationGas(op *userop.UserOpInput, strategy *model.Strategy) (*big.Int, error) {
	//op.SetPreVerificationGas(types.DUMMAY_PREVERIFICATIONGAS_BIGINT)
	//op.SetSignature(types.DUMMY_SIGNATURE_BYTE)
	//Simulate the `packUserOp(p)` function and return a byte slice.
	opValue := *op
	_, userOPPack, err := opValue.PackUserOpForMock(strategy.GetStrategyEntryPointVersion())
	if err != nil {
		return nil, err
	}
	//Calculate the length of the packed byte sequence and convert it to the number of characters.
	lengthInWord := math.Ceil(float64(len(userOPPack)) / 32)
	var callDataConst float64
	for _, b := range userOPPack {
		if b == byte(0) {
			callDataConst += types.GasOverHand.ZeroByte
		} else {
			callDataConst += types.GasOverHand.NonZeroByte
		}
	}
	floatRes := math.Round(callDataConst + types.GasOverHand.Fixed/types.GasOverHand.BundleSize + types.GasOverHand.PerUserOp + types.GasOverHand.PerUserOpWord*lengthInWord)
	floatVal := new(big.Float).SetFloat64(floatRes)
	result := new(big.Int)
	floatVal.Int(result)
	return result, err
}

func getInputData(op *userop.UserOpInput) ([]byte, error) {
	return nil, nil
}
