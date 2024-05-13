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
	//1.if validate Paymaster
	//TODO implement me
	return nil
}

type StarkNetPayService struct {
}

func (s StarkNetPayService) Pay() error {
	//TODO implement me
	return nil
}

func (s StarkNetPayService) RecordAccount() {
	//TODO implement me
	return
}

func (s StarkNetPayService) getReceipt() {
	//TODO implement me
	return
}
