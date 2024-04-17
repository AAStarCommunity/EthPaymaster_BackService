package network

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
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
func getBasicPreVerificationGas(op userop.BaseUserOp) *big.Int {
	op.SetPreVerificationGas(types.DUMMAY_PREVERIFICATIONGAS_BIGINT)
	op.SetSignature(types.DUMMY_SIGNATURE_BYTE)
	op.PackUserOpForMock()
	return big.NewInt(0)
}
