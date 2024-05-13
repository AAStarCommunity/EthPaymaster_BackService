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
}

type Path string

const (
	Auth      Path = "api/auth"
	Healthz   Path = "api/healthz"
	Paymaster Path = "api/v1/paymaster"
)
