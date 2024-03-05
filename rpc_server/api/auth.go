package api

import (
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/middlewares"
	"github.com/gin-gonic/gin"
)

// Auth
// @Tags Auth
// @Description Get AccessToken By ApiKey
// @Accept json
// @Product json
// @Param credential body model.ClientCredential true "AccessToken Model"
// @Router /api/auth [post]
// @Success 200
func Auth(ctx *gin.Context) {
	middlewares.GinJwtMiddleware().LoginHandler(ctx)
}
