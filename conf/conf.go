package conf

import (
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"sync"
)

var once sync.Once

type Conf struct {
	Jwt JWT
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

// getConfiguration
func getConfiguration(filePath *string) *Conf {
	if file, err := os.ReadFile(*filePath); err != nil {
		return mappingEnvToConf(nil)
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

	return conf
}
