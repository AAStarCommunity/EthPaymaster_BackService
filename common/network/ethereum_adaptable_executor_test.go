package network

import (
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_common/contract/contract_entrypoint_v06"
	"AAStarCommunity/EthPaymaster_BackService/common/paymaster_data"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/sirupsen/logrus"
	"math/big"
	"testing"
)

func TestEthereumAdaptableExecutor(t *testing.T) {
	conf.BasicStrategyInit("../../conf/basic_strategy_dev_config.json")
	conf.BusinessConfigInit("../../conf/business_dev_config.json")
	logrus.SetLevel(logrus.DebugLevel)
	op, err := user_op.NewUserOp(utils.GenerateMockUservOperation())
	if err != nil {
		t.Error(err)
		return
	}
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"TestEthereumSepoliaClientConnect",
			func(t *testing.T) {
				testEthereumExecutorClientConnect(t, types.EthereumSepolia)
			},
		},

		{
			"TestScrollSepoliaClientConnect",
			func(t *testing.T) {
				testEthereumExecutorClientConnect(t, types.ScrollSepolia)
			},
		},
		{
			"TestOptimismSepoliaClientConnect",
			func(t *testing.T) {
				testEthereumExecutorClientConnect(t, types.OptimismSepolia)
			},
		},
		{
			"TestArbitrumSpeoliaClientConnect",
			func(t *testing.T) {
				testEthereumExecutorClientConnect(t, types.ArbitrumSpeolia)
			},
		},
		{
			"TestGetUseOpHash",
			func(t *testing.T) {
				testGetUserOpHash(t, types.EthereumSepolia, op)
			},
		},
		{
			"TestSepoliaSimulateV06HandleOp",
			func(t *testing.T) {
				testSimulateV06HandleOp(t, types.EthereumSepolia)
			},
		},
		{
			"TestGetPaymasterAndData",
			func(t *testing.T) {
				testGetPaymasterData(t, types.EthereumSepolia, op)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
func testGetUserOpHash(t *testing.T, chain types.Network, input *user_op.UserOpInput) {
	executor := GetEthereumExecutor(chain)
	if executor == nil {
		t.Error("executor is nil")
	}
	strategy := conf.GetBasicStrategyConfig("Ethereum_Sepolia_v06_verifyPaymaster")
	t.Logf("paymaster Address %s", strategy.GetPaymasterAddress())

	res, _, err := executor.GetUserOpHash(input, strategy)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("userOpHash: %s", hex.EncodeToString(res))
}

func testGetPaymasterData(t *testing.T, chain types.Network, input *user_op.UserOpInput) {
	executor := GetEthereumExecutor(chain)
	if executor == nil {
		t.Error("executor is nil")
	}
	strategy := conf.GetBasicStrategyConfig("Ethereum_Sepolia_v06_verifyPaymaster")
	t.Logf("entryPoint Address %s", strategy.GetEntryPointAddress())
	dataInput := paymaster_data.NewPaymasterDataInput(strategy)
	paymasterData, err := executor.GetPaymasterData(input, strategy, dataInput)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("paymasterData: %v", hex.EncodeToString(paymasterData))

}
func testSimulateV06HandleOp(t *testing.T, chain types.Network) {
	sepoliaExector := GetEthereumExecutor(chain)
	op, newErr := user_op.NewUserOp(utils.GenerateMockUservOperation())
	if newErr != nil {
		t.Error(newErr)
		return
	}
	strategy := conf.GetBasicStrategyConfig("Ethereum_Sepolia_v06_verifyPaymaster")
	t.Logf("entryPoint Address %s", strategy.GetEntryPointAddress())
	simulataResult, err := sepoliaExector.SimulateV06HandleOp(*op, strategy.GetEntryPointAddress())
	if err != nil {
		t.Error(err)
		return
	}
	if simulataResult == nil {
		t.Error("simulataResult is nil")
		return
	}
	t.Logf("simulateResult: %v", simulataResult)
}

func TestSimulate(t *testing.T) {
	conf.BasicStrategyInit("../../conf/basic_strategy_dev_config.json")
	conf.BusinessConfigInit("../../conf/business_dev_config.json")
	op, _ := user_op.NewUserOp(utils.GenerateMockUservOperation())
	abi, _ := contract_entrypoint_v06.ContractMetaData.GetAbi()
	var targetAddress common.Address = common.HexToAddress("0x")
	callData, err := abi.Pack("simulateHandleOp", op, targetAddress, []byte{})
	if err != nil {
		t.Error(err)
		return
	}
	executor := GetEthereumExecutor(types.EthereumSepolia)
	//gClient := executor.GethClient
	//client := executor.Client
	gethClinet := executor.GethClient
	entrypoint := common.HexToAddress("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789")

	//res, err := client.CallContract(context.Background(), ethereum.CallMsg{
	//	From: types.DummyAddress,
	//	To:   &entrypoint,
	//	Data: callData,
	//	Gas:  28072,
	//}, nil)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}

	mapAcc := map[common.Address]gethclient.OverrideAccount{
		entrypoint: {
			Code: EntryPointV06DeployCode,
		}, types.DummyAddress: {
			Nonce:   1,
			Balance: big.NewInt(38312000000001),
		},
	}
	t.Logf("dummyAddress %s", types.DummyAddress.String())
	//addre := common.HexToAddress("0xDf7093eF81fa23415bb703A685c6331584D30177")

	res, err := gethClinet.CallContract(context.Background(), ethereum.CallMsg{
		To:   &entrypoint,
		Data: callData,
	}, nil, &mapAcc)
	resStr := hex.EncodeToString(res)
	t.Logf("simulate result: %v", resStr)

	if err != nil {
		t.Error(err)
		return
	}
	//var hex hexutil.Bytes
	//req := utils.EthCallReq{
	//	From: testAddr,
	//	To:   entrypoint,
	//	Data: callData,
	//}
	//
	//a := struct {
	//	Tracer         string                                        `json:"tracer"`
	//	StateOverrides map[common.Address]gethclient.OverrideAccount `json:"stateOverrides"`
	//}{
	//	StateOverrides: mapAcc,
	//}
	//err = client.Client().CallContext(context.Background(), &hex, "debug_traceCall", &req, "latest", a)
	////t.Logf("simulate result: %v", res)
	//
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
}

func testEthereumExecutorClientConnect(t *testing.T, chain types.Network) {
	executor := GetEthereumExecutor(chain)
	if executor == nil {
		t.Error("executor is nil")
	}
	client := executor.Client
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		t.Error(err)
	}
	if chainId == nil {
		t.Error("chainId is nil")
	}
	if chainId.String() != executor.ChainId.String() {
		t.Errorf(" %s chainId not equal %s", chainId.String(), executor.ChainId.String())
	}
	t.Logf("network %s chainId: %s", chain, chainId.String())
}
