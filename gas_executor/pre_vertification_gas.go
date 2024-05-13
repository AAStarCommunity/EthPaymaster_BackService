package gas_executor

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/network/arbitrum"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"golang.org/x/xerrors"
	"math"
	"math/big"
)

var preVerificationGasFuncMap = map[global_const.NewWorkStack]PreVerificationGasFunc{}

type PreVerificationGasFunc = func(preVerificationEstimateInput *model.PreVerificationGasEstimateInput) (*big.Int, error)

func init() {
	preVerificationGasFuncMap[global_const.ArbStack] = ArbitrumPreVerificationGasFunc()
	preVerificationGasFuncMap[global_const.DefaultStack] = DefaultPreVerificationGasFunc()
	preVerificationGasFuncMap[global_const.OpStack] = OPStackPreVerificationGasFunc()
}
func GetPreVerificationGasFunc(stack global_const.NewWorkStack) (PreVerificationGasFunc, error) {
	function, ok := preVerificationGasFuncMap[stack]
	if !ok {
		return nil, xerrors.Errorf("stack %s not support", stack)
	}
	return function, nil
}

// ArbitrumPreVerificationGasFunc
// https://medium.com/offchainlabs/understanding-arbitrum-2-dimensional-fees-fd1d582596c9.
// https://docs.arbitrum.io/build-decentralized-apps/nodeinterface/reference
func ArbitrumPreVerificationGasFunc() PreVerificationGasFunc {
	return func(preVerificationEstimateInput *model.PreVerificationGasEstimateInput) (*big.Int, error) {
		strategy := preVerificationEstimateInput.Strategy
		base, err := getBasicPreVerificationGas(preVerificationEstimateInput.Op, strategy)
		if err != nil {
			return nil, err
		}
		executor := network.GetEthereumExecutor(strategy.GetNewWork())
		estimateOutPut, err := arbitrum.GetArbEstimateOutPut(executor.Client, preVerificationEstimateInput)
		if err != nil {
			return nil, err
		}

		return big.NewInt(0).Add(base, big.NewInt(int64(estimateOutPut.GasEstimateForL1))), nil
	}
}
func DefaultPreVerificationGasFunc() PreVerificationGasFunc {
	return func(preVerificationEstimateInput *model.PreVerificationGasEstimateInput) (*big.Int, error) {
		return getBasicPreVerificationGas(preVerificationEstimateInput.Op, preVerificationEstimateInput.Strategy)
	}
}

// OPStackPreVerificationGasFunc
// The L1 data fee is paid based on the current Ethereum gas price as tracked within the GasPriceOracle smart contract. This gas price is updated automatically by the OP Mainnet protocol.
// https://docs.optimism.io/builders/app-developers/transactions/estimates#execution-gas-fee
// https://docs.optimism.io/stack/transactions/fees#the-l1-data-fee
func OPStackPreVerificationGasFunc() PreVerificationGasFunc {
	return func(preVerificationEstimateInput *model.PreVerificationGasEstimateInput) (*big.Int, error) {
		op := preVerificationEstimateInput.Op
		strategy := preVerificationEstimateInput.Strategy
		gasFeeResult := preVerificationEstimateInput.GasFeeResult
		basicGas, err := getBasicPreVerificationGas(op, strategy)
		if err != nil {
			return nil, err
		}
		executor := network.GetEthereumExecutor(strategy.GetNewWork())
		_, data, err := op.PackUserOpForMock(strategy.GetStrategyEntrypointVersion())
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
func getBasicPreVerificationGas(op *user_op.UserOpInput, strategy *model.Strategy) (*big.Int, error) {
	//op.SetPreVerificationGas(global_const.DUMMAY_PREVERIFICATIONGAS_BIGINT)
	//op.SetSignature(global_const.DUMMY_SIGNATURE_BYTE)
	//Simulate the `packUserOp(p)` function and return a byte slice.
	opValue := *op
	_, userOPPack, err := opValue.PackUserOpForMock(strategy.GetStrategyEntrypointVersion())
	if err != nil {
		return nil, err
	}
	//Calculate the length of the packed byte sequence and convert it to the number of characters.
	lengthInWord := math.Ceil(float64(len(userOPPack)) / 32)
	var callDataConst float64
	for _, b := range userOPPack {
		if b == byte(0) {
			callDataConst += global_const.GasOverHand.ZeroByte
		} else {
			callDataConst += global_const.GasOverHand.NonZeroByte
		}
	}
	floatRes := math.Round(callDataConst + global_const.GasOverHand.Fixed/global_const.GasOverHand.BundleSize +
		global_const.GasOverHand.PerUserOp + global_const.GasOverHand.PerUserOpWord*lengthInWord)
	floatVal := new(big.Float).SetFloat64(floatRes)
	result := new(big.Int)
	floatVal.Int(result)
	return result, err
}
