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
				strategy := conf.GetBasicStrategyConfig("Ethereum_Sepolia_v06_verifyPaymaster")
				testSimulateHandleOp(t, global_const.EthereumSepolia, strategy)
			},
		},
		//{
		//	"TestSepoliaSimulateV07HandleOp",
		//	func(t *testing.T) {
		//		strategy := conf.GetBasicStrategyConfig("Ethereum_Sepolia_v07_verifyPaymaster")
		//		testSimulateHandleOp(t, global_const.EthereumSepolia, strategy)
		//	},
		//},
		{
			"TestGetPaymasterAndData",
			func(t *testing.T) {
				testGetPaymasterData(t, global_const.EthereumSepolia, op)
			},
		},
		{
			"TestEthExecutorGetPrice",
			func(t *testing.T) {
				testGetPrice(t, global_const.EthereumSepolia)
			},
		},
		{
			"TestSepoliaGetUserTokenBalance",
			func(t *testing.T) {
				testGetBalance(t, global_const.EthereumSepolia, userAddresss)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
func testGetBalance(t *testing.T, chain global_const.Network, address common.Address) {
	executor := GetEthereumExecutor(chain)
	if executor == nil {
		t.Error("executor is nil")
	}
	balance, err := executor.GetUserTokenBalance(address, global_const.USDC)
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

func testGetPaymasterData(t *testing.T, chain global_const.Network, input *user_op.UserOpInput) {
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
	version := strategy.GetStrategyEntryPointVersion()
	var simulataResult *model.SimulateHandleOpResult
	if version == global_const.EntrypointV06 {
		simulataResult, err = sepoliaExector.SimulateV06HandleOp(*op, strategy.GetEntryPointAddress())
	} else if version == global_const.EntryPointV07 {
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
