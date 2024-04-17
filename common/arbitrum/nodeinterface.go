package arbitrum

import (
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_common/paymaster_abi"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

//https://docs.arbitrum.io/build-decentralized-apps/nodeinterface/reference

//

var (
	// GasEstimateL1ComponentMethod https://github.com/OffchainLabs/nitro/blob/v2.2.5/nodeInterface/NodeInterface.go#L473C1-L473C47
	GasEstimateL1ComponentMethod = abi.NewMethod(
		"gasEstimateL1Component",
		"gasEstimateL1Component",
		abi.Function,
		"",
		false,
		true,
		abi.Arguments{
			{Name: "to", Type: paymaster_abi.AddressType},
			{Name: "contractCreation", Type: paymaster_abi.BooleanType},
			{Name: "data", Type: paymaster_abi.BytesType},
		},
		abi.Arguments{
			{Name: "gasEstimateForL1", Type: paymaster_abi.Uint64Type},
			{Name: "baseFee", Type: paymaster_abi.Uint256Type},
			{Name: "l1BaseFeeEstimate", Type: paymaster_abi.Uint256Type},
		},
	)

	PrecompileAddress = common.HexToAddress("0x00000000000000000000000000000000000000C8")
)

type GasEstimateL1ComponentOutput struct {
	GasEstimateForL1  uint64
	BaseFee           *big.Int
	L1BaseFeeEstimate *big.Int
}

func GetEstimateL1ComponentMethod(clinet *ethclient.Client) (*GasEstimateL1ComponentOutput, error) {
	return &GasEstimateL1ComponentOutput{}, nil
}
