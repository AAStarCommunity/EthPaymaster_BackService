package router

import (
	"github.com/gin-gonic/gin"
)

func sponsorUserOperation(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
