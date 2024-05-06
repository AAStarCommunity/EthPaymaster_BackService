package routers

import (
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/api"
	v1 "AAStarCommunity/EthPaymaster_BackService/rpc_server/api/v1"
)

var PrivateRouterMaps []RouterMap
var PublicRouterMaps []RouterMap

func init() {
	PrivateRouterMaps = make([]RouterMap, 0)

	PrivateRouterMaps = append(PrivateRouterMaps, RouterMap{string(TryPayUserOperation), []RestfulMethod{POST}, v1.TryPayUserOperation})
	PrivateRouterMaps = append(PrivateRouterMaps, RouterMap{string(GetSupportStrategy), []RestfulMethod{GET}, v1.GetSupportStrategy})
	PrivateRouterMaps = append(PrivateRouterMaps, RouterMap{string(GetSupportEntrypoint), []RestfulMethod{GET}, v1.GetSupportEntrypoint})
	PrivateRouterMaps = append(PrivateRouterMaps, RouterMap{string(EstimateUserOpGas), []RestfulMethod{POST}, v1.EstimateUserOpGas})
	PrivateRouterMaps = append(PrivateRouterMaps, RouterMap{string(Paymaster), []RestfulMethod{POST}, v1.Paymaster})
	PublicRouterMaps = append(PublicRouterMaps, RouterMap{string(Auth), []RestfulMethod{POST}, api.Auth})
	PublicRouterMaps = append(PublicRouterMaps, RouterMap{string(Healthz), []RestfulMethod{GET, HEAD, OPTIONS}, api.Healthz})
}

type Path string

const (
	TryPayUserOperation  Path = "api/v1/try-pay-user-operation"
	GetSupportStrategy   Path = "api/v1/get-support-strategy"
	GetSupportEntrypoint Path = "api/v1/get-support-entrypoint"
	EstimateUserOpGas    Path = "api/v1/estimate-user-operation-gas"
	Auth                 Path = "api/auth"
	Healthz              Path = "api/healthz"
	Paymaster            Path = "api/v1/paymaster"
)
