package gas_service

import (
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"AAStarCommunity/EthPaymaster_BackService/common/userop"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComputeGas(t *testing.T) {
	userOp, newErr := userop.NewUserOp(utils.GenerateMockUserOperation(), types.EntrypointV06)
	assert.NoError(t, newErr)
	strategy := dashboard_service.GetStrategyById("1")
	gas, err := ComputeGas(userOp, strategy)
	assert.NoError(t, err)
	assert.NotNil(t, gas)
	jsonBypte, _ := json.Marshal(gas)
	fmt.Println(string(jsonBypte))
}
