package promise

import "io"

type Promise interface {
	Desc(arguments []Constant) string
	Eval(arguments []Constant, logger *Logger) bool
}

type Argument interface {
	GetValue(arguments []Constant) string
	String() string
}

type Logger struct {
	Stdout io.Writer
	Stderr io.Writer
}
