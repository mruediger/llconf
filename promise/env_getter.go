package promise

import (
	"os"
)

type EnvGetter struct {
	Name string
}

func (envGetter EnvGetter) GetValue(arguments []Constant, vars *Variables) string {
	value := os.Getenv(envGetter.Name)
	return value
}

func (envGetter EnvGetter) String() string {
	return "env->$" + envGetter.Name + "("+ os.Getenv(envGetter.Name) + ")"
}
