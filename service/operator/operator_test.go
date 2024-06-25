package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/network"
	"AAStarCommunity/EthPaymaster_BackService/common/user_op"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"AAStarCommunity/EthPaymaster_BackService/envirment"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"AAStarCommunity/EthPaymaster_BackService/sponsor_manager"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestOperator(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.InitConfig("../../config/basic_strategy_config.json", "../../config/basic_config.json", "../../config/secret_config.json")
	logrus.SetLevel(logrus.DebugLevel)
	immutableRequest := getMockTryPayUserOpRequest()
	os.Setenv("Env", "unit")
	mockRequestNotSupport1559 := getMockTryPayUserOpRequest()
	mockRequestNotSupport1559.UserOp["maxPriorityFeePerGas"] = mockRequestNotSupport1559.UserOp["maxFeePerGas"]
	sponsor_manager.Init()
	dashboard_service.Init()
	envirment.Environment.SetUnitEnv()
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"testStrategyGenerate",
			func(t *testing.T) {
				testStrategyGenerate(t, immutableRequest)
			},
		},
		{
			"testEstimateUserOpGas",
			func(t *testing.T) {
				testGetEstimateUserOpGas(t, immutableRequest)
			},
		},
		{
			"testGetSupportEntrypointExecute",
			func(t *testing.T) {
				testGetSupportEntrypointExecute(t)
			},
		},
		{
			"Test_ScrollSepoliaV06Verify_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequestNotSupport1559.StrategyCode = string(global_const.StrategyCodeScrollSepoliaV06Verify)
				mockRequest := getMockTryPayUserOpRequest()

				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_EthereumSepoliaV06Verify_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest := getMockTryPayUserOpRequest()

				mockRequest.StrategyCode = string(global_const.StrategyCodeEthereumSepoliaV06Verify)

				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_OptimismSepoliaV06Verify_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest := getMockTryPayUserOpRequest()

				mockRequest.StrategyCode = string(global_const.StrategyCodeOptimismSepoliaV06Verify)

				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_ArbitrumSpeoliaV06Verify_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest := getMockTryPayUserOpRequest()

				mockRequest.StrategyCode = string(global_const.StrategyCodeArbitrumSepoliaV06Verify)

				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_BaseSepoliaV06Verify_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest := getMockTryPayUserOpRequest()
				mockRequest.StrategyCode = string(global_const.StrategyCodeArbitrumSepoliaV06Verify)

				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_EthereumSepoliaV06Erc20_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest := getMockTryPayUserOpRequest()

				mockRequest.UserPayErc20Token = global_const.TokenTypeUSDT
				mockRequest.StrategyCode = string(global_const.StrategyCodeEthereumSepoliaV06Erc20)
				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_OpSepoliaV06Erc20_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest := getMockTryPayUserOpRequest()

				mockRequest.UserPayErc20Token = global_const.TokenTypeUSDT
				mockRequest.StrategyCode = string(global_const.StrategyCodeOptimismSepoliaV06Erc20)
				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_ArbSepoliaV06Erc20_TryPayUserOpExecute",
			func(t *testing.T) {
				mockRequest := getMockTryPayUserOpRequest()
				mockRequest.UserPayErc20Token = global_const.TokenTypeUSDT
				mockRequest.StrategyCode = string(global_const.StrategyCodeArbitrumSepoliaV06Erc20)
				testTryPayUserOpExecute(t, mockRequest)
			},
		},
		{
			"Test_NoSpectCode_TryPayUserOpExecute",
			func(t *testing.T) {
				request := model.UserOpRequest{
					StrategyCode: "123__GhhSA",
					Network:      global_const.EthereumSepolia,
					UserOp:       *utils.GenerateMockUservOperation(),
				}
				testTryPayUserOpExecute(t, &request)
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, tt.test)
	}

}
func testGetEstimateUserOpGas(t *testing.T, request *model.UserOpRequest) {
	result, err := GetEstimateUserOpGas(request)
	if err != nil {
		t.Error(err)
		return
	}
	resultJson, _ := json.Marshal(result)
	fmt.Printf("Result: %v", string(resultJson))
}
func testStrategyGenerate(t *testing.T, request *model.UserOpRequest) {
	strategy, err := StrategyGenerate(request)
	if err != nil {
		t.Error(err)
		return
	}
	strategyJson, _ := json.Marshal(strategy)
	fmt.Printf("Strategy: %v", string(strategyJson))
}
func testGetSupportEntrypointExecute(t *testing.T) {
	res, err := GetSupportEntrypointExecute("ethereum-sepolia")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}
func testTryPayUserOpExecute(t *testing.T, request *model.UserOpRequest) {
	result, err := TryPayUserOpExecute(&model.ApiKeyModel{
		UserId:                        5,
		ProjectSponsorPaymasterEnable: true,
		UserPayPaymasterEnable:        true,
		Erc20PaymasterEnable:          true,
	}, request)
	if err != nil {
		t.Fatal(err)
		return
	}
	resultJson, _ := json.Marshal(result)
	t.Logf("Result: %v", string(resultJson))

	executor := network.GetEthereumExecutor(result.NetWork)
	if executor == nil {
		t.Error("executor is nil")
		return
	}
	userOp, err := user_op.NewUserOp(&request.UserOp)
	paymasterDataStr := result.UserOpResponse.PayMasterAndData

	paymasterData, err := utils.DecodeStringWithPrefix(paymasterDataStr)
	userOp.PaymasterAndData = paymasterData
	if err != nil {
		t.Error(err)
		return
	}
	if result.EntrypointVersion == global_const.EntrypointV07 {
		//TODO
	} else {
		userOp.VerificationGasLimit = utils.ConvertHexToBigInt(result.UserOpResponse.VerificationGasLimit)
		userOp.PreVerificationGas = utils.ConvertHexToBigInt(result.UserOpResponse.PreVerificationGas)
		userOp.MaxFeePerGas = utils.ConvertHexToBigInt(result.UserOpResponse.MaxFeePerGas)
		userOp.MaxPriorityFeePerGas = utils.ConvertHexToBigInt(result.UserOpResponse.MaxPriorityFeePerGas)
		userOp.CallGasLimit = utils.ConvertHexToBigInt(result.UserOpResponse.CallGasLimit)
		address := common.HexToAddress(result.EntryPointAddress)
		jsonUserOP, err := json.Marshal(userOp)
		if err != nil {
			t.Error(err)
			return
		}
		t.Logf("jsonUserOP: %v", string(jsonUserOP))
		result, err := executor.SimulateV06HandleOp(userOp, &address)
		if err != nil {
			t.Fatal(err)
			return
		}
		resultJson, _ := json.Marshal(result)
		t.Logf("Result: %v", string(resultJson))
	}
}

func getMockTryPayUserOpRequest() *model.UserOpRequest {
	return &model.UserOpRequest{
		StrategyCode: "123__GhhSA",
		Network:      global_const.EthereumSepolia,
		UserOp:       *utils.GenerateMockUservOperation(),
	}
}

func TestWSclient(t *testing.T) {
	os.Setenv("Env", "unit")

	t.Logf("Env: %v", os.Getenv("Env"))
	//TODO
	url := "wss://eth-sepolia.g.alchemy.com/v2/wKeLycGxgYRykgf0aGfcpEkUtqyLQg4v"
	wsClient, err := ethclient.Dial(url)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(wsClient)
}
