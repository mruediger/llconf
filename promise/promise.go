package promise

import (
	"io"
	"os"
)

type Promise interface {
	Desc(arguments []Constant) string
	Eval(arguments []Constant, ctx *Context) bool
	New(children []Promise, args []Argument) (Promise,error)
}

type Argument interface {
	GetValue(arguments []Constant, vars *Variables) string
	String() string
}

type Context struct {
	Logger Logger
	Vars   Variables
	Args   []string
	InDir  string
}

type Logger struct {
	Stdout io.Writer
	Stderr io.Writer
	Info   io.Writer
	Changes []ExecType
	Tests []ExecType
}

func NewContext() Context {
	return Context{
		Logger : Logger{
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Info: os.Stdout,
			Changes: []ExecType{},
			Tests: []ExecType{},
		},
		Vars : make(map[string]string),
		InDir : "",
	}
}
