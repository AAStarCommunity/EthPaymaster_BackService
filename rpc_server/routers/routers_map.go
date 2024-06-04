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
	PublicRouterMaps = append(PublicRouterMaps, RouterMap{string(Auth), []RestfulMethod{POST}, api.Auth})
	PublicRouterMaps = append(PublicRouterMaps, RouterMap{string(Healthz), []RestfulMethod{GET, HEAD, OPTIONS}, api.Healthz})
	PublicRouterMaps = append(PublicRouterMaps, RouterMap{string(GetSponsorLog), []RestfulMethod{GET}, v1.GetSponsorDepositAndWithdrawTransactions})
	PublicRouterMaps = append(PublicRouterMaps, RouterMap{string(DepositSponsor), []RestfulMethod{POST}, v1.DepositSponsor})
	PublicRouterMaps = append(PublicRouterMaps, RouterMap{string(WithdrawSponsor), []RestfulMethod{POST}, v1.WithdrawSponsor})
	PublicRouterMaps = append(PublicRouterMaps, RouterMap{string(GetSponsorData), []RestfulMethod{GET}, v1.GetSponsorMetaData})
}

type Path string

const (
	Auth            Path = "api/auth"
	Healthz         Path = "api/healthz"
	Paymaster       Path = "api/v1/paymaster/:network"
	GetSponsorLog   Path = "api/v1/paymaster_sponsor/deposit_log"
	DepositSponsor  Path = "api/v1/paymaster_sponsor/deposit"
	WithdrawSponsor Path = "api/v1/paymaster_sponsor/withdraw"
	GetSponsorData  Path = "api/v1/paymaster_sponsor/data"
)
