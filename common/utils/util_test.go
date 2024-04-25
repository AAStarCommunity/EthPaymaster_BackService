package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
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

//func TestEthereumSign(t *testing.T) {
//	messageStr := "hello world"
//	messageByte := []byte(messageStr)
//	privateKey, _ := crypto.GenerateKey()
//	publicKey := privateKey.Public()
//	publicKeyBytes := crypto.FromECDSAPub(publicKey.(*ecdsa.PublicKey))
//
//	sign, err := GetSign(messageByte, privateKey)
//	if err != nil {
//		t.Errorf("has Error %s", err)
//		return
//	}
//	t.Logf("sign: %x\n", sign)
//
//	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), sign)
//	if err != nil {
//		t.Errorf("has Error %s", err)
//		return
//	}
//	matches := bytes.Equal(sigPublicKey, publicKeyBytes)
//	t.Logf("matches: %t\n", matches)
//}
