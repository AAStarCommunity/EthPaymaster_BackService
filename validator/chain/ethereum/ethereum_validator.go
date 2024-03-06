package ethereum

import (
	"AAStarCommunity/EthPaymaster_BackService/validator/chain"
)

type Validator struct {
	*chain.Base
}

func (e *Validator) IsSupport() bool {
	//TODO implement me
	panic("implement me")
}

func (e *Validator) PreValidate() (err error) {
	//TODO implement me
	panic("implement me")
}

func (e *Validator) AfterGasValidate() (err error) {
	//TODO implement me
	panic("implement me")
}
