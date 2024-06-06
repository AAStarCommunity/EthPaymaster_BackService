package middlewares

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// GenericRecoveryHandler represents the generic error(panic) process
func GenericRecoveryHandler() gin.HandlerFunc {
	DefaultErrorWriter := &PanicExceptionRecord{}
	return gin.RecoveryWithWriter(DefaultErrorWriter, func(c *gin.Context, err interface{}) {
		errInfo := fmt.Sprintf("[panic]: err : [%v] , stack :[%v]", err, utils.GetCurrentGoroutineStack())
		logrus.Error(errInfo)
		logrus.Errorf("%v", errInfo)
		returnError := fmt.Sprintf("%v", err)
		model.GetResponse().SetHttpCode(http.StatusInternalServerError).FailCode(c, http.StatusInternalServerError, returnError)
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
