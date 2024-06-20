package routers

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/common/utils"
	"AAStarCommunity/EthPaymaster_BackService/docs"
	"AAStarCommunity/EthPaymaster_BackService/envirment"
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

var hexParam validator.Func = func(fl validator.FieldLevel) bool {
	param, ok := fl.Field().Interface().(string)
	if ok {
		return utils.ValidateHex(param)
	}
	return true
}

// SetRouters set routers
func SetRouters() (routers *gin.Engine) {
	routers = gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("hexParam", hexParam)
		if err != nil {
			return nil
		}
	}
	buildMod(routers)
	buildRoute(routers)
	routers.NoRoute(func(ctx *gin.Context) {
		model.GetResponse().SetHttpCode(http.StatusNotFound).FailCode(ctx, http.StatusNotFound)
	})

	return routers
}

// buildRouters build routers Init API AND Middleware
func buildRoute(routers *gin.Engine) {
	// build http routers and middleware
	routers.Use(middlewares.GenericRecoveryHandler())

	routers.Use(middlewares.LogHandler())

	routers.Use(middlewares.PvMetrics())
	routers.Use(middlewares.CorsHandler())
	//build the routers not need api access like auth or Traffic limit
	buildRouters(routers, PublicRouterMaps)

	routers.Use(middlewares.ApiVerificationHandler())
	buildRouters(routers, PrivateRouterMaps)
}

// buildMod set Mode by envirment
func buildMod(routers *gin.Engine) {
	buildSwagger(routers)
	// prod mode
	if envirment.Environment.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		//gin.DefaultWriter = io.Discard // disable gin log
		return
	}

	// dev mod
	if envirment.Environment.IsDevelopment() {
		gin.SetMode(gin.DebugMode)
		return
	}
}

// buildSwagger build swagger
func buildSwagger(router *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
