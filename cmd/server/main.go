package main

import (
	"AAStarCommunity/EthPaymaster_BackService/conf"
	"AAStarCommunity/EthPaymaster_BackService/envirment"
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/routers"
	"flag"
	"fmt"
	"os"
	"strings"
)

var aPort = flag.String("port", "", "Port")

// runMode running mode
// @string: Port
func runMode() string {
	flag.Parse()

	if len(*aPort) == 0 {
		*aPort = os.Getenv("port")
	}

	if len(*aPort) == 0 {
		*aPort = ":80"
	}

	if !strings.HasPrefix(*aPort, ":") {
		*aPort = ":" + *aPort
	}

	return *aPort
}

// @contact.name   AAStar Support
// @contact.url    https://aastar.xyz
// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
// @description Type 'Bearer \<TOKEN\>' to correctly set the AccessToken
// @BasePath /api
func main() {
	Init()
	port := runMode()
	_ = routers.SetRouters().Run(port)
}
func Init() {
	strategyPath := fmt.Sprintf("./conf/basic_strategy_%s_config.json", strings.ToLower(envirment.Environment.Name))
	conf.BasicStrategyInit(strategyPath)
	businessConfigPath := fmt.Sprintf("./conf/business_%s_config.json", strings.ToLower(envirment.Environment.Name))
	conf.BusinessConfigInit(businessConfigPath)
}
