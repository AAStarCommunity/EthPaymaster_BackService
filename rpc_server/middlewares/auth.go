package middlewares

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ApiKey struct {
	Key string `form:"apiKey" json:"apiKey" binding:"required"`
}

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Query("apiKey")
		if apiKey == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "ApiKey Is Required Go To Dashboard (dashboard.aastar.io) To Get It"})
			c.Abort()
			return
		}
		apiModel, err := dashboard_service.GetAPiInfoByApiKey(apiKey)
		if err != nil {
			logrus.Errorf("GetAPiInfoByApiKey err: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can Not Find Your Api Key"})
			c.Abort()
			return
		}
		if apiModel.Disable {
			c.JSON(http.StatusForbidden, gin.H{"error": "Api Key Is Disabled"})
			c.Abort()
			return
		}

		c.Set(global_const.ContextKeyApiMoDel, apiModel)
	}
}
