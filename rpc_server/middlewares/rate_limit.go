package middlewares

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/api/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

const (
	DefaultLimit rate.Limit = 50 // limit `DefaultLimit` requests per second
	DefaultBurst int        = 50 // burst size, for surge traffic
)

var limiter map[string]*rate.Limiter

// RateLimiterByApiKeyHandler represents the rate limit by each ApiKey for each api calling
func RateLimiterByApiKeyHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if exists, current := utils.CurrentUser(ctx); exists {
			apiKeyModel, _ := ctx.Get(global_const.ContextKeyApiMoDel)
			defaultLimit := DefaultLimit
			if apiKeyModel != nil {
				defaultLimit = apiKeyModel.(*model.ApiKeyModel).RateLimit
			}
			if limiting(&current, defaultLimit) {
				ctx.Next()
			} else {
				_ = ctx.AbortWithError(http.StatusTooManyRequests, errors.New("too many requests"))
			}
		} else {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("401 Unauthorized"))
		}
	}
}
func clearLimiter(apiKey *string) {
	delete(limiter, *apiKey)
}

func limiting(apiKey *string, defaultLimit rate.Limit) bool {

	var l *rate.Limiter
	if limit, ok := limiter[*apiKey]; ok {
		l = limit
	} else {
		l = rate.NewLimiter(defaultLimit, DefaultBurst)
		limiter[*apiKey] = l
	}

	return l.Allow()
}

func init() {
	limiter = make(map[string]*rate.Limiter, 100)
}
