package promise

import (
	"os"
  "log"
	"bytes"
)

type Promise interface {
	Desc(arguments []Constant) string
	Eval(arguments []Constant, ctx *Context, stack string) bool
	New(children []Promise, args []Argument) (Promise, error)
}

type Argument interface {
	GetValue(arguments []Constant, vars *Variables) string
	String() string
}

type Context struct {
	Logger *Logger
	ExecOutput *bytes.Buffer
	Vars   Variables
	Args   []string
	Env    []string
	InDir  string
	Debug  bool
}

func NewContext() Context {
	return Context{
		Logger:  &Logger{
			Info: log.New(os.Stdout, "llconf (info)", log.LstdFlags),
			Error: log.New(os.Stderr, "llconf (err)", log.LstdFlags|log.Lshortfile),
			Changes: 0,
			Tests:  0,
		},
		Vars:  make(map[string]string),
		InDir: "",
	}
}

type Logger struct {
	Info    *log.Logger
	Error   *log.Logger
	Changes int
	Tests   int
}

func NewLogger() *Logger {
	return &Logger{
		Info: log.New(os.Stdout, "llconf (info)", log.LstdFlags),
		Error: log.New(os.Stderr, "llconf (err)", log.LstdFlags|log.Lshortfile),
		Changes: 0,
		Tests:   0,
	}
}
