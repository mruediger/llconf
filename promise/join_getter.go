package promise

import (
	"strings"
)

type JoinArgument struct {
	Args []Argument
}

func (this JoinArgument) GetValue(arguments []Constant, vars *Variables) string {
	result := ""
	for _, arg := range this.Args {
		result += arg.GetValue(arguments, vars)
	}

	return result
}

func (this JoinArgument) String() string {
	args := []string{}
	for _, arg := range this.Args {
		args = append(args, arg.String())
	}
	return "joinargs-> " + strings.Join(args, " + ")
}
