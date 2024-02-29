package middlewares

import (
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// GenericRecovery represents the generic error(panic) process
func GenericRecovery() gin.HandlerFunc {
	DefaultErrorWriter := &PanicExceptionRecord{}
	return gin.RecoveryWithWriter(DefaultErrorWriter, func(c *gin.Context, err interface{}) {
		errStr := ""
		if conf.Environment.Debugger {
			errStr = fmt.Sprintf("%v", err)
		}
		models.GetResponse().SetHttpCode(http.StatusInternalServerError).FailCode(c, http.StatusInternalServerError, errStr)
	})
}

// PanicExceptionRecord represents the record of panic exception
type PanicExceptionRecord struct{}

func (p *PanicExceptionRecord) Write(b []byte) (n int, err error) {
	s1 := "An error occurred in the server's internal code："
	var build strings.Builder
	build.WriteString(s1)
	build.Write(b)
	errStr := build.String()
	return len(errStr), errors.New(errStr)
}
