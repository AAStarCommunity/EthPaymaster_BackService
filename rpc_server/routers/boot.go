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

	// middlewares
	handlers := make([]gin.HandlerFunc, 0)
	handlers = append(handlers, middlewares.GenericRecovery())
	if conf.Environment.IsDevelopment() {
		handlers = append(handlers, gin.Logger())
	}
	handlers = append(handlers, middlewares.CorsHandler())

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
	routers.Use(handlers...)

	buildRouters(routers) // build http routers

	routers.NoRoute(func(ctx *gin.Context) {
		models.GetResponse().SetHttpCode(http.StatusNotFound).FailCode(ctx, http.StatusNotFound)
	})

	return
}

// buildSwagger build swagger
func buildSwagger(router *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
