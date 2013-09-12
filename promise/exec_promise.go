package promise

import (
	"os"
	"os/exec"
	"strings"
	"fmt"
)

type ExecType int

const (
	ExecChange ExecType = iota
	ExecTest
)

func (t ExecType) Name() string {
	switch t {
	case ExecChange:
		return "change"
	case ExecTest:
		return "test"
	default:
		return "unknown"
	}
}

func (t ExecType) ReportResult(logger *Logger, result bool) {
	fmt.Println(t.Name())
	if t == ExecChange {
		fmt.Println("added change")
		logger.Changes = append(logger.Changes, ExecType(3))
	}
}

type ExecPromise struct {
	Type ExecType
	Arguments []Argument
}

func (p ExecPromise) getCommand(arguments []Constant, vars *Variables) *exec.Cmd {
	largs := p.Arguments
	dir,largs := largs[0].GetValue(arguments, vars), largs[1:]
	
	cmd := ""
	filestat,error := os.Stat(dir)
	if error != nil || !filestat.IsDir() {
		cmd = dir
		dir = os.Getenv("PWD")
	} else {
		cmd, largs = largs[0].GetValue(arguments, vars), largs[1:]
	}

	args := []string{}
	
	for _,argument := range(largs) {
		args = append(args,argument.GetValue(arguments, vars))
	}

	command := exec.Command(cmd,args...)
	command.Dir = dir
	return command
}

func (p ExecPromise) Desc(arguments []Constant) string {
	largs := p.Arguments
	dir,largs := largs[0].String(), largs[1:]
	cmd := ""
		
	filestat,error := os.Stat(dir)
	if error != nil || !filestat.IsDir() {
		cmd = dir
		dir = os.Getenv("PWD")
	} else {
		cmd, largs = largs[0].String(), largs[1:]
	}

	args := make([]string, len(largs))
	for i,v := range largs {
		args[i] = v.String()
	}
	
	return "(" + p.Type.Name() + " in_dir(" + dir + ") <" + cmd + " [" + strings.Join(args,", ") + "] >)"
}

func (p ExecPromise) Eval(arguments []Constant, logger *Logger, vars *Variables) bool {
	command := p.getCommand(arguments, vars)
	command.Stdout = logger.Stdout
	command.Stderr = logger.Stderr

	logger.Stdout.Write([]byte(strings.Join(command.Args, " ") + "\n"))
	
	err := command.Run()

	result := (err == nil)
	p.Type.ReportResult(logger, result)
	return result;
}

type PipePromise struct {
	Execs []ExecPromise
}

func (p PipePromise) Desc(arguments []Constant) string {
	retval := "(pipe"
	for _,v := range(p.Execs) {
		retval += " " + v.Desc(arguments)
	}
	return retval + ")"
}

func (p PipePromise) Eval(arguments []Constant, logger *Logger, vars *Variables) bool {
	commands := []*exec.Cmd{}
	cstrings := []string{}
	
	for _,v := range(p.Execs) {
		cmd :=  v.getCommand(arguments, vars)
		cstrings = append(cstrings, strings.Join(cmd.Args, " "))
		commands = append(commands, cmd)
	}

	for i, command := range(commands[:len(commands) - 1]) {
		out, err := command.StdoutPipe()
		if err != nil {
			logger.Stdout.Write([]byte(err.Error()))
			return false
		}
		command.Start()
		commands[i + 1].Stdin = out
	}

	logger.Stdout.Write([]byte(strings.Join(cstrings, " | ") + "\n"))
	
	last_cmd := commands[len(commands) - 1]
	last_cmd.Stdout = logger.Stdout
	last_cmd.Stderr = logger.Stderr

	err := last_cmd.Run()

	for _, command := range(commands[:len(commands) - 1]) {
		command.Wait()
	}
	
	if (err != nil) {
		return false
	} else {
		return true
	}
}
