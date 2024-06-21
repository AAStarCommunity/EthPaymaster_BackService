package model

import "golang.org/x/time/rate"

type ApiKeyModel struct {
	Disable                       bool       `json:"disable"`
	ApiKey                        string     `json:"api_key"`
	RateLimit                     rate.Limit `json:"rate_limit"`
	UserId                        int64      `json:"user_id"`
	NetWorkLimitEnable            bool       `json:"network_limit_enable"`
	DomainWhitelist               []string   `json:"domain_whitelist"`
	IPWhiteList                   []string   `json:"ip_white_list"`
	PaymasterEnable               bool       `json:"paymaster_enable"`
	Erc20PaymasterEnable          bool       `json:"erc20_paymaster_enable"`
	ProjectSponsorPaymasterEnable bool       `json:"project_sponsor_paymaster_enable"`
	UserPayPaymasterEnable        bool       `json:"user_pay_paymaster_enable"`
}
