package paymaster_data_generator

import "AAStarCommunity/EthPaymaster_BackService/common/model"

type Erc20PaymasterGenerator struct {
}

func (e Erc20PaymasterGenerator) GeneratePayMaster(strategy *model.Strategy, userOp *model.UserOperation) (string, error) {
	//ERC20:[0-1]pay type，[1-21]paymaster address，[21-53]token Amount
	return "0x", nil
}
