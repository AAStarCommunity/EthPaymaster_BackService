package middlewares

import (
	"github.com/gin-gonic/gin"
)

// LogHandler log handler
func LogHandler() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/api/healthz"},
	})
}
