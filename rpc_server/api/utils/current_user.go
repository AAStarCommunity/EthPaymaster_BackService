package utils

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// CurrentUser is a util tool for getting current user(ApiKey) from each rpc request
func CurrentUser(ctx *gin.Context) (exists bool, user string) {

	defer func() {
		if r := recover(); r != nil {
			exists = false
		}
	}()

	mapping := ctx.MustGet("JWT_PAYLOAD").(jwt.MapClaims)

	user = mapping["jti"].(string)

	exists = true

	return
}
