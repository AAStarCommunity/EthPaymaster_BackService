package middlewares

import "github.com/gin-gonic/gin"
import "github.com/gin-contrib/cors"

// CorsHandler 跨域处理中间件
func CorsHandler() gin.HandlerFunc {
	return cors.Default()
}
