package conf

import (
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"sync"
)

var once sync.Once

type Conf struct {
	// TODO: Add Conf Structure Here
}

var conf *Conf

// getConf 读取配置
// 默认从配置文件取，如果配置文件中的db节点内容为空，则从环境变量取
// 如果配置文件不存在，则db从环境变量取，其他值使用默认值
func getConf() *Conf {
	once.Do(func() {
		if conf == nil {
			filePath := getConfFilePath()
			conf = getConfiguration(filePath)
		}
	})
	return conf
}

// getConfiguration 读取配置
// 从配置文件读取，如果环境变量存在对应值，则取环境变量值
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
