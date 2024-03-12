package paymaster_data_generator

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/types"
	"golang.org/x/xerrors"
)

type VerifyingPaymasterGenerator struct {
}

func (v VerifyingPaymasterGenerator) GeneratePayMaster(strategy *model.Strategy, userOp *model.UserOperation, gasResponse *model.ComputeGasResponse, extra map[string]any) ([]byte, error) {
	//verifying:[0-1]pay type，[1-21]paymaster address，[21-85]valid timestamp，[85-] signature
	signature, ok := extra["signature"]
	if !ok {
		return nil, xerrors.Errorf("signature not found")
	}
	res := "0x" + string(types.PayTypeVerifying) + strategy.PayMasterAddress + "" + signature.(string)
	return []byte(res), nil
}
