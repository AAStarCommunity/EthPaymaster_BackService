package api

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"github.com/gin-gonic/gin"
	"time"
)

// Healthz
// @Tags Healthz
// @Description Get Healthz
// @Accept json
// @Product json
// @Param
// @Router /api/health [post,get,put,delete]
// @Success 200
func Healthz(c *gin.Context) {
	response := model.GetResponse()
	response.WithDataSuccess(c, gin.H{
		"hello":       "Eth Paymaster",
		"environment": conf.Environment.Name,
		"time":        time.Now(),
		"version":     "v1.0.0",
	})
	response.Success(c)
}
