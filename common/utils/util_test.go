package utils

import (
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestGenerateKeypair(t *testing.T) {
	privateKey, _ := crypto.GenerateKey()
	privateKeyHex := crypto.FromECDSA(privateKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Printf("privateKeyHex: %x\n", privateKeyHex)
	fmt.Printf("publicKey: %x\n", publicKeyBytes)
	fmt.Printf("address: %s\n", address)
}
func TestSignUserOp(t *testing.T) {
	//privateKeyHex: 1d8a58126e87e53edc7b24d58d1328230641de8c4242c135492bf5560e0ff421
	//publicKey: 044eaed6b1f16e60354156fa334a094affc76d7b7061875a0b04290af9a14cc14ce2bce6ceba941856bd55c63f8199f408fff6495ce9d4c76899055972d23bdb3e
	//address: 0x0E1375d18a4A2A867bEfe908E87322ad031386a6
	_, newErr := userop.NewUserOp(GenerateMockUserOperation())
	if newErr != nil {
		fmt.Println(newErr)
	}

}
func TestNewUserOp(t *testing.T) {
	userOp, newErr := userop.NewUserOp(GenerateMockUserOperation())
	if newErr != nil {
		fmt.Println(newErr)
	}
	res, err := json.Marshal(userOp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("res: %s\n", res)
}
func TestToEthSignedMessageHash(t *testing.T) {
	strByte, err := hex.DecodeString("4bd85fb8854a6bd9dfb18cf88a5bba4daf9bc65f4b8ac00a706f426d40498302")
	if err != nil {
		fmt.Printf("has Error %s", err)
		return
	}
	afterStrByte := ToEthSignedMessageHash(strByte)
	fmt.Printf("afterStrByte: %x\n", afterStrByte)
	afterStr := hex.EncodeToString(afterStrByte)
	fmt.Printf("afterStr: %s\n", afterStr)
}

func TestValidate(t *testing.T) {
	//userOp := GenerateMockUserOperation()
	//assert.True(t, ValidateHex(userOp.Sender))
}
