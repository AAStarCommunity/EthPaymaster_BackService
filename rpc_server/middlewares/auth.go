package middlewares

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiKey struct {
	Key string `form:"apiKey" json:"apiKey" binding:"required"`
}

func ApiVerificationHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Query("apiKey")
		if apiKey == "" {
			_ = c.AbortWithError(http.StatusForbidden, errors.New("ApiKey is mandatory, visit to https://dashboard.aastar.io for more detail"))
			return
		}
		apiModel, err := dashboard_service.GetAPiInfoByApiKey(apiKey)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("can Not Find Your Api Key"))
			return
		}
		if apiModel.Disable {
			_ = c.AbortWithError(http.StatusForbidden, errors.New("api Key Is Disabled"))
			return
		}
		if !apiModel.PaymasterEnable {
			_ = c.AbortWithError(http.StatusForbidden, errors.New("api Key Is Disabled Paymaster"))
			return
		}
		if !VerifyRateLimit(*apiModel) {
			_ = c.AbortWithError(http.StatusTooManyRequests, errors.New("too many requests"))
			return
		}

		if apiModel.IPWhiteList != nil && apiModel.IPWhiteList.Cardinality() > 0 {
			clientIp := c.ClientIP()
			if !apiModel.IPWhiteList.Contains(clientIp) {
				_ = c.AbortWithError(http.StatusForbidden, errors.New("ip not in whitelist"))
				return
			}
		}
		if apiModel.DomainWhitelist != nil && apiModel.DomainWhitelist.Cardinality() > 0 {
			domain := c.Request.Host
			if !apiModel.DomainWhitelist.Contains(domain) {
				_ = c.AbortWithError(http.StatusForbidden, errors.New("domain not in whitelist"))
				return
			}
		}
		c.Set(global_const.ContextKeyApiMoDel, apiModel)
	}
}
