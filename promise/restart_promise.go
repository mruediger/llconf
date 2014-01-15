package promise

import (
	"io"
	"os"
	"errors"
	"os/exec"
	"strings"
	"bitbucket.org/kardianos/osext"
)

type RestartPromise struct {
	NewExe Argument
}

func (p RestartPromise) New(children []Promise, args []Argument) (Promise,error) {
	if len(args) != 1 {
		return nil,errors.New("(restart) needs exactly 1 argument")
	}

	return RestartPromise{args[0]},nil
}

func (p RestartPromise) Desc( arguments []Constant ) string {
	args := make([]string, len(arguments))
	for i,v := range arguments {
		args[i] = v.String()
	}
	return "(restart " + strings.Join(args, ", ") + ")"
}


func (p RestartPromise) Eval( arguments []Constant, ctx *Context) bool {
	newexe := p.NewExe.GetValue(arguments, &ctx.Vars)
	if _, err := os.Stat(newexe); err != nil {
		ctx.Logger.Stderr.Write([]byte(err.Error()))
		return false
	}

	exe,err := osext.Executable()
	if err != nil {
		return false
	}

	os.Rename(newexe, exe)
	ctx.Logger.Stdout.Write([]byte("restartet llconf"))
	p.restartLLConf(exe, ctx.Args, os.Stdout, os.Stderr)
	os.Exit(0)
	return true;
}

func (p RestartPromise) restartLLConf(exe string, args []string, stdout, stderr io.Writer) bool {
	cmd := exec.Command(exe, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Start()
	if err != nil {
		stderr.Write([]byte(err.Error()))
		return false
	} else {
		return true
	}
}
