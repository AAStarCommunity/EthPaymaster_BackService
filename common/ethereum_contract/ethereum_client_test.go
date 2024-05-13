package ethereum_contract

import (
	"AAStarCommunity/EthPaymaster_BackService/common/ethereum_contract/contract/paymaster_verifying_erc20_v07"
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"testing"
)

func TestPaymasterV07(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config.InitConfig("../../config/basic_strategy_config.json", "../../config/basic_config.json", "../../config/secret_config.json")
	network := config.GetNewWorkClientURl(global_const.EthereumSepolia)
	contractAddress := common.HexToAddress("0x3Da96267B98a33267249734FD8FFeC75093D3085")
	t.Logf("network URL %s", network)
	client, err := ethclient.Dial(network)
	if err != nil {
		t.Fatalf("Error: %v", err)
		return
	}
	id, er := client.ChainID(context.Background())
	if er != nil {
		t.Fatalf("Error: %v", er)
		return
	}
	t.Log(id)
	contractInstance, err := paymaster_verifying_erc20_v07.NewContract(contractAddress, client)
	if err != nil {
		t.Fatalf("Error: %v", err)
		return
	}
	//writeInstance, err := paymaster_verifying_erc20_v07.NewContractTransactor(contractAddress, client)
	//if err != nil {
	//	t.Fatalf("Error: %v", err)
	//	return
	//}
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			"testV07Deposit",
			func(t *testing.T) {
				testV07Deposit(t, contractInstance)
			},
		},

		{
			"parsePaymaster",
			func(t *testing.T) {
				parsePaymasterData(t, contractInstance)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func parsePaymasterData(t *testing.T, contractInstance *paymaster_verifying_erc20_v07.Contract) {
	paymasterData := "3da96267b98a33267249734fd8ffec75093d308500000000004c4b40000000000000000000000000001e84800000000000000000000000000000000000000000000000000000000000000000000000006c7bacd00000000000000000000000000000000000000000000000000000000065ed355000000000000000000000000086af7fa0d8b0b7f757ed6cdd0e2aadb33b03be580000000000000000000000000000000000000000000000000000000000000000293df680d08a6d4da0bb7c0ba6d65af835b31f727e83b30e470a697c886597a50e96c2db45aa54b5f83c977745af6b948e86fbabf0fa96f5670e382b7586ac121b"
	paymasterDataByte, er := hex.DecodeString(paymasterData)
	if er != nil {
		t.Fatalf("Error: %v", er)
		return
	}
	res, err := contractInstance.ParsePaymasterAndData(&bind.CallOpts{}, paymasterDataByte)
	if err != nil {
		t.Fatalf("Error: %v", err)
		return
	}
	resJson, _ := json.Marshal(res)
	t.Log(string(resJson))

}

func testV07Deposit(t *testing.T, contractInstance *paymaster_verifying_erc20_v07.Contract) {
	deopsit, err := contractInstance.GetDeposit(&bind.CallOpts{})
	if err != nil {
		t.Fatalf("Error: %v", err)
		return
	}
	t.Log(deopsit)

	verifier, err := contractInstance.Verifier(&bind.CallOpts{})
	if err != nil {
		t.Fatalf("Error: %v", err)
		return
	}
	t.Log(verifier)
}
