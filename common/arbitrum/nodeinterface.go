package arbitrum

import (
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_common/paymaster_abi"
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/xerrors"
	"math/big"
)

//https://docs.arbitrum.io/build-decentralized-apps/nodeinterface/reference

//https://medium.com/offchainlabs/understanding-arbitrum-2-dimensional-fees-fd1d582596c9

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

func GetArbEstimateOutPut(client *ethclient.Client, preVerificationEstimateInput *network.PreVerificationEstimateInput) (*GasEstimateL1ComponentOutput, error) {
	strategy := preVerificationEstimateInput.Strategy
	simulteaOpCallData := preVerificationEstimateInput.SimulateOpResult.SimulateUserOpCallData
	methodIputeData, err := GasEstimateL1ComponentMethod.Inputs.Pack(strategy.GetEntryPointAddress(), false, append(simulteaOpCallData))

	if err != nil {
		return nil, err
	}
	req := map[string]any{
		"from": global_const.EmptyAddress,
		"to":   PrecompileAddress,
		"data": hexutil.Encode(append(GasEstimateL1ComponentMethod.ID, methodIputeData...)),
	}
	var outPut any
	if err := client.Client().Call(&outPut, "eth_call", &req, "latest"); err != nil {
		return nil, err
	}
	outPutStr, ok := outPut.(string)
	if !ok {
		return nil, xerrors.Errorf("gasEstimateL1Component: cannot assert type: hex is not of type string")
	}
	data, err := hexutil.Decode(outPutStr)
	if err != nil {
		return nil, xerrors.Errorf("gasEstimateL1Component: %s", err)
	}
	outputArgs, err := GasEstimateL1ComponentMethod.Outputs.Unpack(data)
	if err != nil {
		return nil, xerrors.Errorf("gasEstimateL1Component: %s", err)
	}

	return &GasEstimateL1ComponentOutput{
		GasEstimateForL1:  outputArgs[0].(uint64),
		BaseFee:           outputArgs[1].(*big.Int),
		L1BaseFeeEstimate: outputArgs[2].(*big.Int),
	}, nil
}
