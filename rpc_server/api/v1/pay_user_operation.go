package v1

import (
	"github.com/gin-gonic/gin"
)

// TryPayUserOperation
// @Tags Sponsor
// @Description sponsor the userOp
// @Accept json
// @Product json
// @Router /api/v1/try-pay-user-operation [post]
// @Success 200
func TryPayUserOperation(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
