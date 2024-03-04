package api

import (
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/models"
	"github.com/gin-gonic/gin"
	"time"
)

// Healthz
// @Tags Healthz
// @Description Get Healthz
// @Router /api/healthz [get]
// @Success 200
func Healthz(c *gin.Context) {
	response := models.GetResponse()
	response.WithDataSuccess(c, gin.H{
		"hello":       "Eth Paymaster",
		"environment": conf.Environment.Name,
		"time":        time.Now(),
		"version":     "v1.0.0",
	})
}
