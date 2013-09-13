package promise

import (
	"io"
	"os"
)

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
	Info   io.Writer
	Changes []ExecType
	Tests []ExecType
}

func NewStdoutLogger() Logger {
	return Logger{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Info: os.Stdout,
		Changes: []ExecType{},
		Tests: []ExecType{},
	}
}

