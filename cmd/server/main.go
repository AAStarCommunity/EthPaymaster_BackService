package main

import (
	"AAStarCommunity/EthPaymaster_BackService/config"
	"AAStarCommunity/EthPaymaster_BackService/envirment"
	"AAStarCommunity/EthPaymaster_BackService/rpc_server/routers"
	"AAStarCommunity/EthPaymaster_BackService/service/dashboard_service"
	"AAStarCommunity/EthPaymaster_BackService/sponsor_manager"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

const (
	strategyPath    = "./config/basic_strategy_config.json"
	basicConfigPath = "./config/basic_config.json"
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
// @BasePath /api
func main() {
	secretPath := os.Getenv("secret_config_path")
	if secretPath == "" {
		secretPath = "./config/secret_config.json"
	}
	initEngine(strategyPath, basicConfigPath, secretPath)
	port := runMode()
	os.Getenv("secret_config_path")
	_ = Engine.Run(port)
}

func initEngine(strategyPath string, basicConfigPath string, secretPath string) {
	logrus.Infof("secretPath: %s", secretPath)
	config.InitConfig(strategyPath, basicConfigPath, secretPath)
	if envirment.Environment.IsDevelopment() {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	dashboard_service.Init()
	sponsor_manager.Init()
	logrus.Infof("Environment: %s", envirment.Environment.Name)
	logrus.Infof("Debugger: %v", envirment.Environment.Debugger)
	Engine = routers.SetRouters()
}
