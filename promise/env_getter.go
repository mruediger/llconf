package promise

import (
	"os"
)

type EnvGetter struct {
	Name string
}

func (envGetter EnvGetter) GetValue(arguments []Constant) string {
	value := os.Getenv(envGetter.Name)
	if value == "" {
		panic("didn't find environment variable " + envGetter.Name)
	}
	return value
}

func (envGetter EnvGetter) String() string {
	return "env->$" + envGetter.Name + "("+ os.Getenv(envGetter.Name) + ")"
}
