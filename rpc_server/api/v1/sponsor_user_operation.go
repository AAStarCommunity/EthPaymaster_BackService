package v1

import (
	"github.com/gin-gonic/gin"
)

// SponsorUserOperation
// @Tags Sponsor
// @Description sponsor the userOp
// @Accept json
// @Product json
// @Router /api/v1/sponsor-user-operation [post]
// @Success 200
func SponsorUserOperation(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
