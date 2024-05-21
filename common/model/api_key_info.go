package model

import "golang.org/x/time/rate"

type ApiKeyModel struct {
	Disable   bool       `gorm:"column:disable;type:bool" json:"disable"`
	ApiKey    string     `gorm:"column:api_key;type:varchar(255)" json:"api_key"`
	RateLimit rate.Limit `gorm:"column:rate_limit;type:int" json:"rate_limit"`
}
