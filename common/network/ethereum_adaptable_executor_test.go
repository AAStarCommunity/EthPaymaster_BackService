package network

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/paymaster_data"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
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
	userAddresss := common.HexToAddress("0xFD44DF0Fe211d5EFDBe1423483Fcb3FDeF84540f")
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"TestEthereumSepoliaClientConnect",
			func(t *testing.T) {
				testEthereumExecutorClientConnect(t, global_const.EthereumSepolia)
			},
		},

		{
			"TestScrollSepoliaClientConnect",
			func(t *testing.T) {
				testEthereumExecutorClientConnect(t, global_const.ScrollSepolia)
			},
		},
		{
			"TestOptimismSepoliaClientConnect",
			func(t *testing.T) {
				testEthereumExecutorClientConnect(t, global_const.OptimismSepolia)
			},
		},
		{
			"TestArbitrumSpeoliaClientConnect",
			func(t *testing.T) {
				testEthereumExecutorClientConnect(t, global_const.ArbitrumSpeolia)
			},
		},
		{
			"TestGetUseOpHash",
			func(t *testing.T) {
				testGetUserOpHash(t, global_const.EthereumSepolia, op)
			},
		},
		{
			"TestSepoliaSimulateV06HandleOp",
			func(t *testing.T) {
				strategy := conf.GetBasicStrategyConfig(global_const.StrategyCodeEthereumSepoliaV06Verify)
				testSimulateHandleOp(t, global_const.EthereumSepolia, strategy)
			},
		},
		{
			"TestSepoliaSimulateV06HandleOp",
			func(t *testing.T) {
				strategy := conf.GetBasicStrategyConfig(global_const.StrategyCodeOptimismSepoliaV06Verify)
				testSimulateHandleOp(t, global_const.OptimismSepolia, strategy)
			},
		},
		//{
		//	"TestScrollSepoliaSimulateV06HandleOp",
		//	func(t *testing.T) {
		//		strategy := conf.GetBasicStrategyConfig(global_const.StrategyCodeScrollSepoliaV06Verify)
		//		testSimulateHandleOp(t, global_const.ScrollSepolia, strategy)
		//	},
		//},
		{
			"TestSepoliaSimulateV07HandleOp",
			func(t *testing.T) {
				strategy := conf.GetBasicStrategyConfig(global_const.StrategyCodeEthereumSepoliaV07Verify)
				testSimulateHandleOp(t, global_const.EthereumSepolia, strategy)
			},
		},
		{
			"TestGetPaymasterAndDataV07",
			func(t *testing.T) {
				strategy := conf.GetBasicStrategyConfig(global_const.StrategyCodeEthereumSepoliaV07Verify)
				testGetPaymasterData(t, global_const.EthereumSepolia, op, strategy)
			},
		},
		{
			"TestEthExecutorGetPrice",
			func(t *testing.T) {
				testGetPrice(t, global_const.EthereumSepolia)
			},
		},
		{
			"TestScrollExecutorGetPrice",
			func(t *testing.T) {
				testGetPrice(t, global_const.ScrollSepolia)
			},
		},
		{
			"TestSepoliaGetUserTokenBalance",
			func(t *testing.T) {
				testGetBalance(t, global_const.EthereumSepolia, userAddresss)
			},
		},
		{
			"checkContractAddressAccess",
			func(t *testing.T) {
				address := common.HexToAddress("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789")
				testCheckContractAddressAccess(t, global_const.EthereumSepolia, address)
			},
		},
		{
			"checkContractAddressAccess",
			func(t *testing.T) {
				address := common.HexToAddress("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789")
				testCheckContractAddressAccess(t, global_const.ScrollSepolia, address)
			},
		},
		{
			"TestEstimateUserOpCallGas",
			func(t *testing.T) {
				strategy := conf.GetBasicStrategyConfig("Ethereum_Sepolia_v06_verifyPaymaster")
				entrypointAddress := strategy.GetEntryPointAddress()
				testEstimateUserOpCallGas(t, global_const.EthereumSepolia, op, entrypointAddress)
			},
		},
		{
			"TestEstimateCreateSenderGas",
			func(t *testing.T) {
				strategy := conf.GetBasicStrategyConfig("Ethereum_Sepolia_v06_verifyPaymaster")
				entrypointAddress := strategy.GetEntryPointAddress()
				testEstimateCreateSenderGas(t, global_const.EthereumSepolia, op, entrypointAddress)
			},
		},
		{
			"TestOptimismGetL1DataFee",
			func(t *testing.T) {
				stategy := conf.GetBasicStrategyConfig("Optimism_Sepolia_v06_verifyPaymaster")

				testGetL1DataFee(t, global_const.OptimismSepolia, *op, stategy.GetStrategyEntrypointVersion())
			},
		},
		{
			"TestOpPreVerificationGasFunc",
			func(t *testing.T) {

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func testGetL1DataFee(t *testing.T, chain global_const.Network, input user_op.UserOpInput, version global_const.EntrypointVersion) {
	executor := GetEthereumExecutor(chain)
	if executor == nil {
		t.Error("executor is nil")
	}
	_, data, err := input.PackUserOpForMock(version)
	if err != nil {
		t.Error(err)
		return
	}
	l1DataFee, err := executor.GetL1DataFee(data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("l1DataFee: %v", l1DataFee)
}
func testEstimateUserOpCallGas(t *testing.T, chain global_const.Network, userOpParam *user_op.UserOpInput, entpointAddress *common.Address) {
	executor := GetEthereumExecutor(chain)
	if executor == nil {
		t.Error("executor is nil")
	}
	gasResult, err := executor.EstimateUserOpCallGas(entpointAddress, userOpParam)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("gasResult: %v", gasResult)
}
func testEstimateCreateSenderGas(t *testing.T, chain global_const.Network, userOpParam *user_op.UserOpInput, entrypointAddress *common.Address) {
	executor := GetEthereumExecutor(chain)
	if executor == nil {
		t.Error("executor is nil")
	}
	gasResult, err := executor.EstimateCreateSenderGas(entrypointAddress, userOpParam)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("gasResult: %v", gasResult)
}
func testCheckContractAddressAccess(t *testing.T, chain global_const.Network, address common.Address) {
	executor := GetEthereumExecutor(chain)
	if executor == nil {
		t.Error("executor is nil")
	}
	res, err := executor.CheckContractAddressAccess(&address)
	if err != nil {
		t.Error(err)
		return
	}
	if !res {
		t.Error("checkContractAddressAccess failed")
	}
}
func testGetBalance(t *testing.T, chain global_const.Network, address common.Address) {
	executor := GetEthereumExecutor(chain)
	if executor == nil {
		t.Error("executor is nil")
	}
	balance, err := executor.GetUserTokenBalance(address, global_const.TokenTypeUSDC)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("balance: %v", balance)
}

func testGetPrice(t *testing.T, chain global_const.Network) {
	executor := GetEthereumExecutor(chain)
	if executor == nil {
		t.Error("executor is nil")
	}
	price, err := executor.GetGasPrice()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("price: %v", price)
}
func testGetUserOpHash(t *testing.T, chain global_const.Network, input *user_op.UserOpInput) {
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

func testGetPaymasterData(t *testing.T, chain global_const.Network, input *user_op.UserOpInput, strategy *model.Strategy) {
	executor := GetEthereumExecutor(chain)
	if executor == nil {
		t.Error("executor is nil")
	}
	t.Logf("entryPoint Address %s", strategy.GetEntryPointAddress())
	dataInput := paymaster_data.NewPaymasterDataInput(strategy)
	dataInput.PaymasterPostOpGasLimit = global_const.DummyPaymasterPostopGaslimitBigint
	dataInput.PaymasterVerificationGasLimit = global_const.DummyPaymasterVerificationgaslimitBigint
	paymasterData, err := executor.GetPaymasterData(input, strategy, dataInput)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("paymasterData: %v", hex.EncodeToString(paymasterData))

}
func testSimulateHandleOp(t *testing.T, chain global_const.Network, strategy *model.Strategy) {
	sepoliaExector := GetEthereumExecutor(chain)
	op, newErr := user_op.NewUserOp(utils.GenerateMockUservOperation())
	if newErr != nil {
		t.Error(newErr)
		return
	}
	dataInput := paymaster_data.NewPaymasterDataInput(strategy)
	paymasterData, err := sepoliaExector.GetPaymasterData(op, strategy, dataInput)
	if err != nil {
		t.Error(err)
		return
	}
	op.PaymasterAndData = paymasterData
	t.Logf("entryPoint Address %s", strategy.GetEntryPointAddress())
	version := strategy.GetStrategyEntrypointVersion()
	t.Logf("version: %s", version)
	var simulataResult *model.SimulateHandleOpResult
	if version == global_const.EntrypointV06 {
		simulataResult, err = sepoliaExector.SimulateV06HandleOp(*op, strategy.GetEntryPointAddress())
	} else if version == global_const.EntrypointV07 {
		simulataResult, err = sepoliaExector.SimulateV07HandleOp(*op, strategy.GetEntryPointAddress())
	}

	if err != nil {
		t.Error(err)
		return
	}
	if simulataResult == nil {
		t.Error("simulataResult is nil")
		return
	}
	t.Logf("simulateResult: %v", simulataResult)
	callData := simulataResult.SimulateUserOpCallData
	t.Logf("callData: %v", hex.EncodeToString(callData))
}

func testEthereumExecutorClientConnect(t *testing.T, chain global_const.Network) {
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

func TestStr(t *testing.T) {

	res, err := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000001e41413331207061796d6173746572206465706f73697420746f6f206c6f770000")

	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("res: %v", res)
	var resStr string

	json.Unmarshal(res, &resStr)
	t.Logf("resStr: %v", resStr)
}
