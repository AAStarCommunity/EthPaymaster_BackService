package envirment

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"os"
	"sync"
)

var once sync.Once
var Environment *model.Env

func init() {
	envName := model.DevEnv
	if len(os.Getenv(model.EnvKey)) > 0 {
		envName = os.Getenv(model.EnvKey)
	}
	Environment = &model.Env{
		Name: envName,
		Debugger: func() bool {
			return envName != model.ProdEnv
		}(),
	}
}
