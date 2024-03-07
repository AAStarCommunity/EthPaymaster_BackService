package utils

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateMockUserOperation() *model.UserOperationItem {
	//TODO use config
	return &model.UserOperationItem{
		Sender:               "0x4A2FD3215420376DA4eD32853C19E4755deeC4D1",
		Nonce:                "1",
		InitCode:             "0x",
		CallData:             "0xb61d27f6000000000000000000000000c206b552ab127608c3f666156c8e03a8471c72df000000000000000000000000000000000000000000000000002386f26fc1000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000000",
		CallGasLimit:         "39837",
		VerificationGasList:  "100000",
		PreVerificationGas:   "44020",
		MaxFeePerGas:         "1743509478",
		MaxPriorityFeePerGas: "1500000000",
		Signature:            "0x760868cd7d9539c6e31c2169c4cab6817beb8247516a90e4301e929011451658623455035b83d38e987ef2e57558695040a25219c39eaa0e31a0ead16a5c925c1c",
	}
}
func GenerateUserOperation() *model.UserOperationItem {
	return &model.UserOperationItem{}
}

func SignUserOp(privateKeyHex string, userOp *model.UserOperationItem) ([]byte, error) {

	serializedUserOp, err := json.Marshal(userOp)
	if err != nil {
		return nil, err
	}

	signature, err := SignMessage(privateKeyHex, string(serializedUserOp))
	if err != nil {
		return nil, err
	}

	return signature, nil
}
func SignMessage(privateKeyHex string, message string) ([]byte, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	// Hash the message
	hash := crypto.Keccak256([]byte(message))

	// Sign the hash
	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return nil, err
	}

	return signature, nil
}
