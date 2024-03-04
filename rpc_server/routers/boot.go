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

	buildMod(routers)
	buildRoute(routers)
	routers.NoRoute(func(ctx *gin.Context) {
		models.GetResponse().SetHttpCode(http.StatusNotFound).FailCode(ctx, http.StatusNotFound)
	})

	return
}
func buildRoute(routers *gin.Engine) {
	// build http routers and middleware
	routers.Use(middlewares.GenericRecoveryHandler())
	if conf.Environment.IsDevelopment() {
		routers.Use(middlewares.LogHandler())
	}
	routers.Use(middlewares.CorsHandler())
	//build the routers not need api access like auth or Traffic limit
	buildRouters(routers, PublicRouterMaps)

	routers.Use(middlewares.AuthHandler())
	routers.Use(middlewares.RateLimiterByApiKeyHandler())
	buildRouters(routers, PrivateRouterMaps)
}

func buildMod(routers *gin.Engine) {

	// prod mode
	if conf.Environment.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard // disable gin log
		return
	}

	// dev mod
	if conf.Environment.IsDevelopment() {
		gin.SetMode(gin.DebugMode)
		buildSwagger(routers)
		return
	}
}

// buildSwagger build swagger
func buildSwagger(router *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
