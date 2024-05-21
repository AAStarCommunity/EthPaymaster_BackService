package model

import "golang.org/x/time/rate"

type ApiKeyModel struct {
	Disable   bool       `json:"disable"`
	ApiKey    string     `json:"api_key"`
	RateLimit rate.Limit `json:"rate_limit"`
	UserId    int64      `json:"user_id"`
}
