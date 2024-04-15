package network

import "AAStarCommunity/EthPaymaster_BackService/common/model"

type StarknetExecutor struct {
	BaseExecutor
}

func init() {

}
func GetStarknetExecutor() *StarknetExecutor {
	return &StarknetExecutor{}
}
func (executor StarknetExecutor) GetGasPrice() (*model.GasPrice, error) {
	return nil, nil
}
