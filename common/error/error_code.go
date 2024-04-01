package error

type ErrorCode string

const (
	ValidateParamError  ErrorCode = "AA2"
	ValidateUserOpError ErrorCode = "AA3"
	ValidateGasError    ErrorCode = "AA4"
)
