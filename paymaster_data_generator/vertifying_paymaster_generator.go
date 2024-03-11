package paymaster_data_generator

import "AAStarCommunity/EthPaymaster_BackService/common/model"

type VerifyingPaymasterGenerator struct {
}

func (v VerifyingPaymasterGenerator) GeneratePayMaster(strategy *model.Strategy, userOp *model.UserOperationItem) (string, error) {
	//verifying:[0-1]pay type，[1-21]paymaster address，[21-85]valid timestamp，[85-] signature
	return "0x", nil

}
