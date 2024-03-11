package types

type ErrorPrefix string

const (
	ValidateParamError  ErrorPrefix = "AA2"
	ValidateUserOpError ErrorPrefix = "AA3"
	ValidateGasError    ErrorPrefix = "AA4"
)
