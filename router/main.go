package router

import (
	"github.com/gin-gonic/gin"
)

var (
	Engine *gin.Engine
)

func InitRouter() {
	if Engine != nil {
		return
	}

	Engine = gin.Default()
	Engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	Engine.POST("/sponsorUserOperation", sponsorUserOperation)
	Engine.POST("/v1/user/relation/bind", validateUserOperation)

}
