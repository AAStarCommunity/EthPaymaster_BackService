package routers

import (
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/api"
	v1 "AAStarCommunity/EthPaymaster_BackService/rpc_server/api/v1"
)

var PrivateRouterMaps []RouterMap
var PublicRouterMaps []RouterMap

func init() {
	PrivateRouterMaps = make([]RouterMap, 0)
	PrivateRouterMaps = append(PrivateRouterMaps, RouterMap{string(Paymaster), []RestfulMethod{POST}, v1.Paymaster})
	PublicRouterMaps = append(PublicRouterMaps, RouterMap{string(Healthz), []RestfulMethod{GET, HEAD, OPTIONS}, api.Healthz})
	PublicRouterMaps = append(PublicRouterMaps, RouterMap{string(DepositSponsor), []RestfulMethod{POST}, v1.DepositSponsor})
	PublicRouterMaps = append(PublicRouterMaps, RouterMap{string(WithdrawSponsor), []RestfulMethod{POST}, v1.WithdrawSponsor})
	PublicRouterMaps = append(PublicRouterMaps, RouterMap{string(GetTokenPrice), []RestfulMethod{GET}, v1.GetTokenPrice})
}

type Path string

const (
	Healthz         Path = "api/healthz"
	Paymaster       Path = "api/v1/paymaster/:network"
	GetSponsorLog   Path = "api/v1/paymaster_sponsor/deposit_log"
	DepositSponsor  Path = "api/v1/paymaster_sponsor/deposit"
	WithdrawSponsor Path = "api/v1/paymaster_sponsor/withdraw"
	GetSponsorData  Path = "api/v1/paymaster_sponsor/data"
	GetTokenPrice   Path = "api/v1/paymaster_sponsor/token_price"
)
