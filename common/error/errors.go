package error

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}
