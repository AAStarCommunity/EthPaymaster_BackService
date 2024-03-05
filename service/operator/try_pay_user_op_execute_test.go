package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"fmt"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestTryPayUserOpExecute(t *testing.T) {
	request := getMockTryPayUserOpRequest()
	result, err := TryPayUserOpExecute(&request)
	assert.NoError(t, err)
	fmt.Printf("Result: %v", result)
}

func getMockTryPayUserOpRequest() model.TryPayUserOpRequest {
	return model.TryPayUserOpRequest{
		ForceStrategyId: "1",
		UserOperation: model.UserOperationItem{
			Sender:               "0x123",
			Nonce:                "",
			CallGasLimit:         "",
			VerificationGasList:  "",
			PerVerificationGas:   "",
			MaxFeePerGas:         "",
			MaxPriorityFeePerGas: "",
		},
	}
}
