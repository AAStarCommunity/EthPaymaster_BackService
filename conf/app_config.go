package conf

import (
	"AAStarCommunity/EthPaymaster_BackService/common/model"
	"fmt"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"strings"
	"sync"
)

var once sync.Once
var Environment *model.Env

type Conf struct {
	Jwt JWT
}

func EnvInit() {
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

var conf *Conf

// getConf read conf from file
func getConf() *Conf {
	once.Do(func() {
		if conf == nil {
			filePath := getConfFilePath()
			conf = getConfiguration(filePath)
		}
	})
	return conf
}
func getConfFilePath() *string {
	path := fmt.Sprintf("conf/appsettings.%s.yaml", strings.ToLower(Environment.Name))
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		path = fmt.Sprintf("conf/appsettings.yaml")
	}
	return &path
}

// getConfiguration
func getConfiguration(filePath *string) *Conf {
	if file, err := os.ReadFile(*filePath); err != nil {
		return mappingEnvToConf(
			&Conf{
				Jwt: JWT{},
			},
		)
	} else {
		c := Conf{}
		err := yaml.Unmarshal(file, &c)
		if err != nil {
			return mappingEnvToConf(&c)
		}

		return &c
	}
}

func mappingEnvToConf(conf *Conf) *Conf {

	// TODO: read from env
	// e.g. if dummy := os.Getenv("dummy"); len(dummy) > 0 {conf.Dummy = dummy}
	if jwt__security := os.Getenv("jwt__security"); len(jwt__security) > 0 {
		conf.Jwt.Security = jwt__security
	}
	if jwt__realm := os.Getenv("jwt__realm"); len(jwt__realm) > 0 {
		conf.Jwt.Security = jwt__realm
	}
	if jwt__idkey := os.Getenv("jwt__idkey"); len(jwt__idkey) > 0 {
		conf.Jwt.Security = jwt__idkey
	}

	return conf
}
