package main

import (
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/routers"
	"flag"
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

func main() {
	port := runMode()
	_ = routers.SetRouters().Run(port)
}
