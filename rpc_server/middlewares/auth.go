package middlewares

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/envirment"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
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
		Realm:       envirment.GetJwtKey().Realm,
		Key:         []byte(envirment.GetJwtKey().Security),
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
			apiModel, err := dashboard_service.GetAPiInfoByApiKey(apiKey.Key)
			if err != nil {
				return "", err
			}
			err = checkAPIKeyAvailable(apiModel)
			if err != nil {
				return "", err
			}
			return apiKey.Key, nil
		},
		//every Request  will be checked
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// always return true unless the permission feature started
			if data == nil {
				logrus.Errorf("Authorizator id is nil")
				return false
			}
			apiKey := data.(string)
			if apiKey == "" {
				logrus.Errorf("Authorizator id is nil")
				return false
			}

			apiModel, err := dashboard_service.GetAPiInfoByApiKey(apiKey)
			if err != nil {
				c.Set("ERROR_REASON", err.Error())
				return false
			}
			if apiModel == nil {
				c.Set("ERROR_REASON", "API Key is not found")
				return false
			}
			err = checkAPIKeyAvailable(apiModel)
			if err != nil {
				c.Set("ERROR_REASON", err.Error())
				return false
			}
			c.Set(global_const.ContextKeyApiMoDel, apiModel)
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			if c.GetString("ERROR_REASON") != "" {
				message = c.GetString("ERROR_REASON")
			}
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
func checkAPIKeyAvailable(apiModel *model.ApiKeyModel) error {
	if apiModel.Disable {
		return xerrors.Errorf("API Key is disabled")
	}
	return nil
}
