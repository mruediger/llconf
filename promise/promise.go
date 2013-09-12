package promise

import "io"

type Promise interface {
	Desc(arguments []Constant) string
	Eval(arguments []Constant, logger *Logger, vars *Variables) bool
}

type Argument interface {
	GetValue(arguments []Constant, vars *Variables) string
	String() string
}

type Logger struct {
	Stdout io.Writer
	Stderr io.Writer
	Changes []ExecType
	Tests []ExecType
}	
