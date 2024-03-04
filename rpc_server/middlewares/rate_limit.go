package middlewares

import (
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

			if limiting(&current) {
				ctx.Next()
			} else {
				_ = ctx.AbortWithError(http.StatusTooManyRequests, errors.New("too many requests"))
			}
		} else {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("401 Unauthorized"))
		}
	}
}

func limiting(apiKey *string) bool {

	var l *rate.Limiter
	if limit, ok := limiter[*apiKey]; ok {
		l = limit
	} else {
		// TODO: different rate config for each current(apiKey) should get from dashboard service
		l = rate.NewLimiter(DefaultLimit, DefaultBurst)
		limiter[*apiKey] = l
	}

	return l.Allow()
}

func init() {
	limiter = make(map[string]*rate.Limiter, 100)
}
