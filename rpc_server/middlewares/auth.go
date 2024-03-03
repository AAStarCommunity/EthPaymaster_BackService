package middlewares

import (
	"AAStarCommunity/EthPaymaster_BackService/conf"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"time"
)

type ApiKey struct {
	Key string `form:"apiKey" json:"apiKey" binding:"required"`
}

var jwtMiddleware *jwt.GinJWTMiddleware

func GinJwtMiddleware() *jwt.GinJWTMiddleware {
	return jwtMiddleware
}

func AuthHandler() gin.HandlerFunc {
	m, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       conf.GetJwtKey().Realm,
		Key:         []byte(conf.GetJwtKey().Security),
		Timeout:     time.Hour * 24,
		MaxRefresh:  time.Hour / 2,
		IdentityKey: "jti",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(string); ok {
				return jwt.MapClaims{
					"jti": v,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var apiKey ApiKey
			if err := c.ShouldBind(&apiKey); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			// TODO: verify if the key is correct
			return apiKey.Key, nil

			// if incorrect
			//return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// always return true unless the permission feature started
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		HTTPStatusMessageFunc: func(e error, c *gin.Context) string {
			return "401 Unauthorized"
		},
	})

	jwtMiddleware = m

	return m.MiddlewareFunc()
}
