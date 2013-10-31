package promise

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"bitbucket.org/kardianos/osext"
)

type RestartPromise struct {
	NewExe Argument
}


func (p RestartPromise) Decs( arguments []Constant ) string {
	args := make([]string, len(arguments))
	for i,v := range arguments {
		args[i] = v.String()
	}
	return "(restart " + strings.Join(args, ", ") + ")"
}


func (p RestartPromise) Eval( arguments []Constant, ctx *Context) bool {

	newexe := p.NewExe.GetValue(arguments, &ctx.Vars)
	if _, err := os.Stat(newexe); err != nil {
		os.Stderr.Write([]byte(err.Error()))
		return false
	}

	exe,err := osext.Executable()
	if err != nil {
		return false
	}

	os.Rename(exe, exe + ".old")
	os.Rename(newexe, exe)

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
