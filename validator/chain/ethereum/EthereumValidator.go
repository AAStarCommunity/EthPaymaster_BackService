package ethereum

import (
	"AAStarCommunity/EthPaymaster_BackService/validator/chain"
)

type EthValidator struct {
	*chain.Base
}

func (e EthValidator) IsSupport() bool {
	//TODO implement me
	panic("implement me")
}

func (e EthValidator) PreValidate() (err error) {
	//TODO implement me
	panic("implement me")
}

func (e EthValidator) AfterGasValidate() (err error) {
	//TODO implement me
	panic("implement me")
}
