package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
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
	signByte, err := SignUserOp("1d8a58126e87e53edc7b24d58d1328230641de8c4242c135492bf5560e0ff421", GenerateMockUserOperation())
	assert.NoError(t, err)
	fmt.Printf("signByte: %x\n", signByte)
	singature := hex.EncodeToString(signByte)
	fmt.Printf("singature: %s\n", singature)

}
