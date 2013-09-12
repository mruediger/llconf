package promise

import (
	"strconv"
)

type ArgGetter struct {
	Position int
}

func (argGetter ArgGetter) GetValue(arguments []Constant, vars *Variables) string {
	if len(arguments) <= argGetter.Position {
		return ""
	}
	return arguments[argGetter.Position].Value
}

func (argGetter ArgGetter) String() string {
	return "arg->" + strconv.Itoa(argGetter.Position)
}
