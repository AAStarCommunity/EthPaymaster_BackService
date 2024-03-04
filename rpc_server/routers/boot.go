package routers

import (
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"AAStarCommunity/EthPaymaster_BackService/docs"
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/middlewares"
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/models"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	"net/http"
)

// SetRouters set routers
func SetRouters() (routers *gin.Engine) {
	routers = gin.New()

	// prod mode
	if conf.Environment.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard // disable gin log
	}

	// dev mod
	if conf.Environment.IsDevelopment() {
		gin.SetMode(gin.DebugMode)
		buildSwagger(routers)
	}

	// use middlewares
	handlers := generateHandlers()
	routers.Use(handlers...)

	// build http routers
	buildRouters(routers)

	routers.NoRoute(func(ctx *gin.Context) {
		models.GetResponse().SetHttpCode(http.StatusNotFound).FailCode(ctx, http.StatusNotFound)
	})

	return
}

func generateHandlers() []gin.HandlerFunc {
	// middlewares
	handlers := make([]gin.HandlerFunc, 0)
	handlers = append(handlers, middlewares.GenericRecoveryHandler())
	if conf.Environment.IsDevelopment() {
		handlers = append(handlers, middlewares.LogHandler())
	}
	handlers = append(handlers, middlewares.CorsHandler())
	handlers = append(handlers, middlewares.AuthHandler())
	handlers = append(handlers, middlewares.RateLimiterByApiKeyByHandler())
	return handlers
}

// buildSwagger build swagger
func buildSwagger(router *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
