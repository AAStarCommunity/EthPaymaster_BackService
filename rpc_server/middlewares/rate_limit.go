package middlewares

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"golang.org/x/time/rate"
)

const (
	DefaultLimit rate.Limit = 50 // limit `DefaultLimit` requests per second
	DefaultBurst int        = 50 // burst size, for surge traffic
)

var limiter map[string]*rate.Limiter

func VerifyRateLimit(keyModel model.ApiKeyModel) bool {
	return limiting(&keyModel.ApiKey, keyModel.RateLimit)
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
