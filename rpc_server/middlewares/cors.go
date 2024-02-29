package middlewares

import "github.com/gin-gonic/gin"
import "github.com/gin-contrib/cors"

// CorsHandler cross domain handler
func CorsHandler() gin.HandlerFunc {
	return cors.Default()
}
