package routers

import (
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/api"
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/middlewares"
	"github.com/gin-gonic/gin"
)

// buildRouters Build Routers
func buildRouters(router *gin.Engine) {

	router.POST("api/auth", api.Auth)

	router.Use(middlewares.AuthHandler())
	{
		router.Use(middlewares.RateLimiterByApiKey())

		for _, routerMap := range RouterMaps {
			for _, method := range routerMap.Methods {
				if method == GET {
					router.GET(routerMap.Url, routerMap.Func)
				} else if method == PUT {
					router.PUT(routerMap.Url, routerMap.Func)
				} else if method == POST {
					router.POST(routerMap.Url, routerMap.Func)
				} else if method == DELETE {
					router.DELETE(routerMap.Url, routerMap.Func)
				} // ignore rest methods
			}
		}
	}
}
