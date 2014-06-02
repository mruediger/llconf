package promise

import (
	"bitbucket.org/kardianos/osext"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type RestartPromise struct {
	NewExe Argument
}

func (p RestartPromise) New(children []Promise, args []Argument) (Promise, error) {
	if len(args) != 1 {
		return nil, errors.New("(restart) needs exactly 1 argument")
	}

	if len(children) != 0 {
		return nil, errors.New("(restart) cannot have nested promises")
	}

	return RestartPromise{args[0]}, nil
}

func (p RestartPromise) Desc(arguments []Constant) string {
	args := make([]string, len(arguments))
	for i, v := range arguments {
		args[i] = v.String()
	}
	return "(restart " + strings.Join(args, ", ") + ")"
}

func (p RestartPromise) Eval(arguments []Constant, ctx *Context, stack string) bool {
	newexe := p.NewExe.GetValue(arguments, &ctx.Vars)
	if _, err := os.Stat(newexe); err != nil {
		ctx.Logger.Error.Print(err.Error())
		return false
	}

	exe, err := osext.Executable()
	if err != nil {
		return false
	}

	os.Rename(newexe, exe)
	ctx.Logger.Info.Print(fmt.Printf("restarted llconf: llconf %v", ctx.Args))

	if _, err := p.restartLLConf(exe, ctx.Args, ctx.ExecOutput, ctx.ExecOutput); err == nil {
		os.Exit(0)
		return true
	} else {
		ctx.Logger.Error.Print(err.Error())
		return false
	}
}
func (p RestartPromise) restartLLConf(exe string, args []string, stdout, stderr io.Writer) (*exec.Cmd, error) {
	cmd := exec.Command(exe, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	err := cmd.Start()
	return cmd, err
}
