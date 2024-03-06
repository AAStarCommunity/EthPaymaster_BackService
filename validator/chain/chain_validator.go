package chain

type Validator interface {
	PreValidate() (err error)
	AfterGasValidate() (err error)
	IsSupport() bool
}
type Base struct {
	name string
}
