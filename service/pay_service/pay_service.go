package pay_service

type PayService interface {
	Pay() error
	RecordAccount()
	getReceipt()
}

type EthereumPayService struct {
}

func (e *EthereumPayService) RecordAccount() {
	//TODO implement me

}

func (e *EthereumPayService) GetReceipt() {
	//TODO implement me

}

func (e *EthereumPayService) Pay() error {
	//TODO implement me
	return nil
}
