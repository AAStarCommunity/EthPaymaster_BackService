package routers

import "github.com/gin-gonic/gin"

type RestfulMethod string

const (
	PUT     RestfulMethod = "PUT"
	GET     RestfulMethod = "GET"
	DELETE  RestfulMethod = "DELETE"
	POST    RestfulMethod = "POST"
	HEAD    RestfulMethod = "HEAD"
	OPTIONS RestfulMethod = "OPTIONS"
)

type RouterMap struct {
	Url     string
	Methods []RestfulMethod
	Func    func(ctx *gin.Context)
}
