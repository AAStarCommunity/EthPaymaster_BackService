package model

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetResponse() *Response {
	return &Response{
		httpCode: http.StatusOK,
		Result: &Result{
			Code:    200,
			Message: "",
			Data:    nil,
			Cost:    "",
		},
	}
}
func BadRequest(ctx *gin.Context, data ...any) {
	GetResponse().withDataAndHttpCode(http.StatusBadRequest, ctx, data)
}

// Success represents response success
func Success(ctx *gin.Context, data ...any) {
	if data != nil {
		GetResponse().WithDataSuccess(ctx, data[0])
		return
	}
	GetResponse().Success(ctx)
}

// Fail represents response failed
func Fail(ctx *gin.Context, code int, message *string, data ...any) {
	var msg string
	if message == nil {
		msg = ""
	} else {
		msg = *message
	}
	if data != nil {
		GetResponse().WithData(data[0]).FailCode(ctx, code, msg)
		return
	}
	GetResponse().FailCode(ctx, code, msg)
}

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Cost    string      `json:"cost"`
}

type Response struct {
	httpCode int
	Result   *Result
}

// Fail represents 5XX error
func (r *Response) Fail(ctx *gin.Context) *Response {
	r.SetCode(http.StatusInternalServerError)
	r.json(ctx)
	return r
}

// FailCode customer error codes
func (r *Response) FailCode(ctx *gin.Context, code int, msg ...string) *Response {
	r.SetCode(code)
	if msg != nil {
		r.WithMessage(msg[0])
	}
	r.json(ctx)
	return r
}

// Success represents the success response
func (r *Response) Success(ctx *gin.Context) *Response {
	r.SetCode(http.StatusOK)
	r.json(ctx)
	return r
}

// WithDataSuccess return success with data
func (r *Response) WithDataSuccess(ctx *gin.Context, data interface{}) *Response {
	r.SetCode(http.StatusOK)
	r.WithData(data)
	r.json(ctx)
	return r
}

func (r *Response) withDataAndHttpCode(code int, ctx *gin.Context, data interface{}) *Response {
	r.SetHttpCode(code)
	if data != nil {
		r.WithData(data)
	}
	r.json(ctx)
	return r
}

// SetCode set business code
func (r *Response) SetCode(code int) *Response {
	r.Result.Code = code
	return r
}

// SetHttpCode set http status code
func (r *Response) SetHttpCode(code int) *Response {
	r.httpCode = code
	return r
}

type defaultRes struct {
	Result any `json:"result"`
}

// WithData represents response with data
func (r *Response) WithData(data interface{}) *Response {
	switch data.(type) {
	case string, int, bool:
		r.Result.Data = &defaultRes{Result: data}
	default:
		r.Result.Data = data
	}
	return r
}

// WithMessage represents returns response with message
func (r *Response) WithMessage(message string) *Response {
	r.Result.Message = message
	return r
}

// json returns HandlerFunc
func (r *Response) json(ctx *gin.Context) {
	r.Result.Cost = time.Since(ctx.GetTime("requestStartTime")).String()
	ctx.AbortWithStatusJSON(r.httpCode, r.Result)
}
