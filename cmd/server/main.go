package main

import (
	"AAStarCommunity/EthPaymaster_BackService/config"
	"AAStarCommunity/EthPaymaster_BackService/envirment"
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/routers"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

var Engine *gin.Engine

// @contact.name   AAStar Support
// @contact.url    https://aastar.xyz
// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
// @description Type 'Bearer \<TOKEN\>' to correctly set the AccessToken
// @BasePath /api
func main() {
	strategyPath := fmt.Sprintf("./config/basic_strategy_%s_config.json", strings.ToLower(envirment.Environment.Name))
	businessConfigPath := fmt.Sprintf("./config/business_%s_config.json", strings.ToLower(envirment.Environment.Name))

	initEngine(strategyPath, businessConfigPath)
	port := runMode()
	_ = Engine.Run(port)
}

func initEngine(strategyPath string, businessConfigPath string) {
	config.BasicStrategyInit(strategyPath)
	config.BusinessConfigInit(businessConfigPath)
	if envirment.Environment.IsDevelopment() {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.Infof("Environment: %s", envirment.Environment.Name)
	logrus.Infof("Debugger: %v", envirment.Environment.Debugger)
	logrus.Infof("Action ENV : [%v]", os.Getenv("GITHUB_ACTIONS"))
	Engine = routers.SetRouters()
}
