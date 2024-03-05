package routers

import (
	"github.com/gin-gonic/gin"
)

// buildRouters Build Routers
func buildRouters(router *gin.Engine, routerMaps []RouterMap) {
	for _, routerMap := range routerMaps {
		for _, method := range routerMap.Methods {
			if method == GET {
				router.GET(routerMap.Url, routerMap.Func)
			} else if method == PUT {
				router.PUT(routerMap.Url, routerMap.Func)
			} else if method == POST {
				router.POST(routerMap.Url, routerMap.Func)
			} else if method == DELETE {
				router.DELETE(routerMap.Url, routerMap.Func)
			} else if method == OPTIONS {
				router.OPTIONS(routerMap.Url, routerMap.Func)
			} else if method == HEAD {
				router.HEAD(routerMap.Url, routerMap.Func)
			}
			// ignore rest methods
		}

	}
}
