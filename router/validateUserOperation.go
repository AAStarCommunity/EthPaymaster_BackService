package router

import "github.com/gin-gonic/gin"

func validateUserOperation(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
