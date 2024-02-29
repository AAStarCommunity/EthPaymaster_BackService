package main

import (
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/routers"
	"flag"
	"os"
	"strings"
)

var aPort = flag.String("port", "", "端口")

// runMode running mode
// @string: 端口
func runMode() string {
	// 优先读取命令行参数，其次使用go env，最后使用默认值
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
