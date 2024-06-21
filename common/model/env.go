package model

import (
	"strings"
)

const EnvKey = "Env"
const ProdEnv = "prod"
const DevEnv = "dev"
const UnitEnv = "unit"

type Env struct {
	Name     string // env Name, like `prod`, `dev` and etc.,
	Debugger bool   // whether to use debugger
}

func (env *Env) IsDevelopment() bool {
	return strings.EqualFold(DevEnv, env.Name) || strings.EqualFold(UnitEnv, env.Name)
}
func (env *Env) IsUnit() bool {
	return strings.EqualFold(UnitEnv, env.Name)
}
func (env *Env) IsProduction() bool {
	return strings.EqualFold(ProdEnv, env.Name)
}

func (env *Env) GetEnvName() *string {
	return &env.Name
}
func (env *Env) SetUnitEnv() {
	env.Name = UnitEnv
}
