package v1

import "github.com/gin-gonic/gin"

// ValidateUserOperation
// @Tags Sponsor
// @Description validate the userOp for sponsor
// @Accept json
// @Product json
// @Router /api/v1/validate-user-operation [post]
// @Success 200
func ValidateUserOperation(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
