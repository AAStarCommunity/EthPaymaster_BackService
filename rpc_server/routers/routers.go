package routers

import v1 "AAStarCommunity/EthPaymaster_BackService/rpc_server/api/v1"

var RouterMaps []RouterMap

func init() {
	RouterMaps = make([]RouterMap, 0)

	RouterMaps = append(RouterMaps, RouterMap{"api/v1/try-pay-user-operation", []RestfulMethod{POST}, v1.TryPayUserOperation})
	RouterMaps = append(RouterMaps, RouterMap{"api/v1/get_support_strategy", []RestfulMethod{POST}, v1.GetSupportStrategy})
	RouterMaps = append(RouterMaps, RouterMap{"api/v1/get_support_entrypoint", []RestfulMethod{POST}, v1.GetSupportEntrypoint})

}
