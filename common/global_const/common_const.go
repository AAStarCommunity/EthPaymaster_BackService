package global_const

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

const (
	DummyPrivateKeyText                  = "0a82406dc7fcf16090e05215ff394c7465608dd1a698632471b1eb37b8ece2f7"
	DummySignature                       = "0x3054659b5e29460a8f3ac9afc3d5fcbe4b76f92aed454b944e9b29e55d80fde807716530b739540e95cfa4880d69f710a9d45910f2951a227675dc1fb0fdf2c71c"
	DummyPaymasterData                   = "0xd93349Ee959d295B115Ee223aF10EF432A8E8523000000000000000000000000000000000000000000000000000000001710044496000000000000000000000000000000000000000000000000000000174158049605bea0bfb8539016420e76749fda407b74d3d35c539927a45000156335643827672fa359ee968d72db12d4b4768e8323cd47443505ab138a525c1f61c6abdac501"
	DummyInitCode                        = "0x9406cc6185a346906296840746125a0e449764545fbfb9cf000000000000000000000000b6bcf9517d193f551d0e3d6860103972dd13de7b0000000000000000000000000000000000000000000000000000000000000000"
	DummyVerificationGas                 = 50000
	DUMMY_PAYMASTER_POSTOP_GASLIMIT      = 2000000
	DUMMY_PAYMASTER_VERIFICATIONGASLIMIT = 5000000
)

var (
	EthWeiFactor                           = new(big.Float).SetInt(big.NewInt(1e18))
	DummySignatureByte                     []byte
	DummyInitCodeByte                      []byte
	DummyReverificationsBigint             = big.NewInt(DummyVerificationGas)
	DummyPaymasterOversimplificationBigint = big.NewInt(DUMMY_PAYMASTER_VERIFICATIONGASLIMIT)
	DummyPaymasterPostoperativelyBigint    = big.NewInt(DUMMY_PAYMASTER_POSTOP_GASLIMIT)
	ThreeBigint                            = big.NewInt(3)
	HundredBigint                          = big.NewInt(100)
	TwoBigint                              = big.NewInt(2)
	HundredPlusOneBigint                   = big.NewInt(110)
	DummyPrivateKey                        *ecdsa.PrivateKey
	DummyAddress                           common.Address
	DummyPaymasterDataByte                 []byte
	DummyMaxFeePerGas                      = big.NewInt(1500012654)
	DummyMaxPriorityFeePerGas              = big.NewInt(1500000000)
	DummyCallGasLimit                      = big.NewInt(21754)
	DummyVerificationGasLimit              = big.NewInt(391733)
	EmptyAddress                           = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

func init() {
	privateKey, err := crypto.HexToECDSA(DummyPrivateKeyText)
	if err != nil {
		panic(err)
	}
	DummyPrivateKey = privateKey
	address := crypto.PubkeyToAddress(DummyPrivateKey.PublicKey)
	DummyAddress = address
	DummyPaymasterDataByte, err = hex.DecodeString(DummyPaymasterData[2:])
	if err != nil {
		panic(err)
	}
}

var GasOverHand = struct {
	//fixed overhead for entire handleOp bundle.
	Fixed float64
	//per userOp overhead, added on top of the above fixed per-bundle
	PerUserOp float64
	//overhead for userOp word (32 bytes) block
	PerUserOpWord float64
	//zero byte cost, for calldata gas cost calculations
	ZeroByte float64
	//non-zero byte cost, for calldata gas cost calculations
	NonZeroByte float64
	//expected bundle size, to split per-bundle overhead between all ops.
	BundleSize float64
	//expected length of the userOp signature.
	sigSize float64
}{
	Fixed:         21000,
	PerUserOp:     18300,
	PerUserOpWord: 4,
	ZeroByte:      4,
	NonZeroByte:   16,
	BundleSize:    1,
	sigSize:       65,
}

func init() {
	signatureByte, err := hex.DecodeString(DummySignature[2:])
	if err != nil {
		panic(err)
	}
	DummySignatureByte = signatureByte
}
