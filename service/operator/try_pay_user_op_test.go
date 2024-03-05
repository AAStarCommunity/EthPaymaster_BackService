package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"

	"testing"
)

func TestTryPayUserOpExecute(t *testing.T) {
	request := getMockTryPayUserOpRequest()
	TryPayUserOpExecute(request)

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
