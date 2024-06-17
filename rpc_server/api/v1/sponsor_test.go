package v1

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/config"
	"AAStarCommunity/EthPaymaster_BackService/sponsor_manager"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestValidateDeposit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
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
	//Only For Test
	//privateKeyHex: 95ecbd0ae2055889d6c40ed042ef846d091313541482e61129f48867cd15fc4a
	//publicKey: 0401e57b7947d19b224a98700dea8bcff3dd556980823aca6182e12eebe64e42ecee690d4c2d6eee6bbf9f82423e38314beec2be0a6833741a917a8abaebf66a76
	//address: 0x3aCF4b1F443a088186Cbd66c5F81479C6e968eCA
	request := &model.DepositSponsorRequest{
		DepositAddress: "0xFD44DF0Fe211d5EFDBe1423483Fcb3FDeF84540f",
		TxHash:         "0x367428ad744c2fd80054283f7143b934d630433f9c40a411d91e65893dbabdf1",
		PayUserId:      "5",
		DepositSource:  "dashboard",
		IsTestNet:      true,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	hash := sha256.New()
	hash.Write(jsonData)
	hashByte := hash.Sum(nil)
	hashHex := hex.EncodeToString(hashByte)
	t.Logf("hash: %v", hashHex)
	signerEoa, err := global_const.NewEoa("95ecbd0ae2055889d6c40ed042ef846d091313541482e61129f48867cd15fc4a")
	if err != nil {
		t.Error(err)
		return
	}
	signatureByte, err := crypto.Sign(accounts.TextHash(hashByte), signerEoa.PrivateKey)
	if err != nil {
		t.Error(err)
		return
	}
	signatureByteHex := hex.EncodeToString(signatureByte)
	t.Logf("signatureByteHex: %v", signatureByteHex)

	err = ValidateSignature(hashHex, signatureByteHex, jsonData, "0x3aCF4b1F443a088186Cbd66c5F81479C6e968eCA")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("ValidateSignature success")
}
func TestDemo(t *testing.T) {
	t.Logf("Demo")
	address := "0xFfDB071C2b58CCC10Ad386f9Bb4E8d3d664CE73c"
	commonAddres := common.HexToAddress(address)
	commonAddres.Hex()
	t.Logf("commonAddres: %v", commonAddres.Hex())

}
