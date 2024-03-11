package paymaster_data_generator

import "AAStarCommunity/EthPaymaster_BackService/common/model"

type PaymasterGenerator interface {
	GeneratePayMaster(strategy *model.Strategy, userOp *model.UserOperationItem) (string, error)
}
