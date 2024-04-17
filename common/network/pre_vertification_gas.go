package network

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"math"
	"math/big"
)

var PreVerificationGasFuncMap = map[types.NewWorkStack]PreVerificationGasFunc{}

type PreVerificationGasFunc = func() (*big.Int, error)

func init() {
	PreVerificationGasFuncMap[types.ARBSTACK] = ArbitrumPreVerificationGasFunc()
	PreVerificationGasFuncMap[types.DEFAULT_STACK] = DefaultPreVerificationGasFunc()
	PreVerificationGasFuncMap[types.OPSTACK] = OPStackPreVerificationGasFunc()
}

// https://medium.com/offchainlabs/understanding-arbitrum-2-dimensional-fees-fd1d582596c9.
func ArbitrumPreVerificationGasFunc() PreVerificationGasFunc {
	return func() (*big.Int, error) {
		return big.NewInt(0), nil
	}
}
func DefaultPreVerificationGasFunc() PreVerificationGasFunc {
	return func() (*big.Int, error) {
		return big.NewInt(0), nil
	}
}
func OPStackPreVerificationGasFunc() PreVerificationGasFunc {
	return func() (*big.Int, error) {
		return big.NewInt(0), nil
	}
}

/**
 * calculate the preVerificationGas of the given UserOperation
 * preVerificationGas (by definition) is the cost overhead that can't be calculated on-chain.
 * it is based on parameters that are defined by the Ethereum protocol for external transactions.
 * @param userOp filled userOp to calculate. The only possible missing fields can be the signature and preVerificationGas itself
 * @param overheads gas overheads to use, to override the default values
 */
func getBasicPreVerificationGas(op userop.BaseUserOp) (*big.Int, error) {
	op.SetPreVerificationGas(types.DUMMAY_PREVERIFICATIONGAS_BIGINT)
	op.SetSignature(types.DUMMY_SIGNATURE_BYTE)
	//Simulate the `packUserOp(p)` function and return a byte slice.
	_, userOPPack, err := op.PackUserOpForMock()
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
