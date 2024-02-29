package routers

import v1 "AAStarCommunity/EthPaymaster_BackService/rpc_server/api/v1"

var RouterMaps []RouterMap

func init() {
	RouterMaps = make([]RouterMap, 0)

	RouterMaps = append(RouterMaps, RouterMap{"api/v1/tryPayUserOperation", []RestfulMethod{POST}, v1.TryPayUserOperation})
}
