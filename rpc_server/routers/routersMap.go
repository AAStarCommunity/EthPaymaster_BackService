package routers

import (
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/api"
	v1 "AAStarCommunity/EthPaymaster_BackService/rpc_server/api/v1"
)

var RouterMaps []RouterMap
var RouterNotAPIAccessMaps []RouterMap

func init() {
	RouterMaps = make([]RouterMap, 0)

	RouterMaps = append(RouterMaps, RouterMap{string(TryPayUserOperation), []RestfulMethod{POST}, v1.TryPayUserOperation})
	RouterMaps = append(RouterMaps, RouterMap{string(GetSupportStrategy), []RestfulMethod{GET}, v1.GetSupportStrategy})
	RouterMaps = append(RouterMaps, RouterMap{string(GetSupportEntrypoint), []RestfulMethod{GET}, v1.GetSupportEntrypoint})
	RouterNotAPIAccessMaps = append(RouterMaps, RouterMap{string(Auth), []RestfulMethod{POST}, api.Auth})
	RouterNotAPIAccessMaps = append(RouterMaps, RouterMap{string(Health), []RestfulMethod{GET}, api.Health})
}

type Path string

const (
	TryPayUserOperation  Path = "api/v1/try-pay-user-operation"
	GetSupportStrategy   Path = "api/v1/get-support-strategy"
	GetSupportEntrypoint Path = "api/v1/get-support-entrypoint"
	Auth                 Path = "api/auth"
	Health               Path = "api/health"
)
