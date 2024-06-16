package middlewares

import (
	"AAStarCommunity/EthPaymaster_BackService/common/global_const"
	"AAStarCommunity/EthPaymaster_BackService/common/model"
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
		apiKeyModelInterface := ctx.MustGet(global_const.ContextKeyApiMoDel)
		defaultLimit := DefaultLimit
		apiKeyModel := apiKeyModelInterface.(*model.ApiKeyModel)
		defaultLimit = apiKeyModel.RateLimit

		if limiting(&apiKeyModel.ApiKey, defaultLimit) {
			ctx.Next()
		} else {
			_ = ctx.AbortWithError(http.StatusTooManyRequests, errors.New("too many requests"))
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
