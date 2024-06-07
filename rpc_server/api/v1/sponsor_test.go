package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"AAStarCommunity/EthPaymaster_BackService/sponsor_manager"
	"testing"
)

func TestValidateDeposit(t *testing.T) {
	config.InitConfig("../../../config/basic_strategy_config.json", "../../../config/basic_config.json", "../../../config/secret_config.json")
	sponsor_manager.Init()
	request := &model.DepositSponsorRequest{
		DepositAddress: "0xFD44DF0Fe211d5EFDBe1423483Fcb3FDeF84540f",
		TxHash:         "0x367428ad744c2fd80054283f7143b934d630433f9c40a411d91e65893dbabdf1",
		PayUserId:      "5",
		DepositSource:  "dashboard",
		IsTestNet:      true,
	}
	sender, amount, err := validateDeposit(request)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("sender: %v, amount: %v", sender.Hex(), amount.String())

}
func TestValidateSignature(t *testing.T) {
	originHash := "0x367428ad744c2fd80054283f7143b934d630433f9c40a411d91e65893dbabdf1"
	signatureHex := "0x367428ad744c2fd80054283f7143b934d630433f9c40a411d91e65893dbabdf1"
	inputJson := []byte("0x367428ad744c2fd80054283f7143b934d630433f9c40a411d91e65893dbabdf1")
	signerAddress := "0xFD44DF0Fe211d5EFDBe1423483Fcb3FDeF84540f"
	err := ValidateSignature(originHash, signatureHex, inputJson, signerAddress)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("ValidateSignature success")
}
