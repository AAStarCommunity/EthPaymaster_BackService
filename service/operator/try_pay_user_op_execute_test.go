package operator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTryPayUserOpExecute(t *testing.T) {
	request := getMockTryPayUserOpRequest()
	result, err := TryPayUserOpExecute(request)
	assert.NoError(t, err)
	resultJson, _ := json.Marshal(result)
	fmt.Printf("Result: %v", string(resultJson))
}

func getMockTryPayUserOpRequest() *model.TryPayUserOpRequest {
	return &model.TryPayUserOpRequest{
		ForceStrategyId: "1",
		UserOp:          *utils.GenerateMockUserOperation(),
	}
}
