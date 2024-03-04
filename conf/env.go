package conf

import (
	"fmt"
	"os"
	"strings"
)

const envKey = "Env"
const ProdEnv = "prod"
const DevEnv = "dev"

type Env struct {
	Name     string // env Name, like `prod`, `dev` and etc.,
	Debugger bool   // whether to use debugger
}

func (env *Env) IsDevelopment() bool {
	return strings.EqualFold("dev", env.Name)
}

func (env *Env) IsProduction() bool {
	return strings.EqualFold("prod", env.Name)
}

func (env *Env) GetEnvName() *string {
	return &env.Name
}

func getConfFilePath() *string {
	path := fmt.Sprintf("conf/appsettings.%s.yaml", strings.ToLower(Environment.Name))
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		path = fmt.Sprintf("conf/appsettings.yaml")
	}
	return &path
}

var Environment *Env

func init() {
	envName := ProdEnv
	if len(os.Getenv(envKey)) > 0 {
		envName = os.Getenv(envKey)
	}
	Environment = &Env{
		Name: envName,
		Debugger: func() bool {
			return envName != ProdEnv
		}(),
	}
}
