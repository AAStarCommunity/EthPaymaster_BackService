package paymaster_data_generator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"encoding/hex"
)

type Erc20PaymasterGenerator struct {
}

func (e *Erc20PaymasterGenerator) GeneratePayMaster(strategy *model.Strategy, userOp *model.UserOperation, gasResponse *model.ComputeGasResponse, extra map[string]any) ([]byte, error) {
	//ERC20:[0-1]pay type，[1-21]paymaster address，[21-53]token Amount
	res := "0x" + string(types.PayTypeERC20) + strategy.PayMasterAddress + gasResponse.TokenCost
	//TODO implement me
	return hex.DecodeString(res)
}
