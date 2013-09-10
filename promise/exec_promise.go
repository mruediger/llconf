package promise

import (
	"os"
	"os/exec"
	"strings"
)

type ExecPromise struct {
	Arguments []Argument
}

func (p ExecPromise) getCommand(arguments []Constant) *exec.Cmd {
	largs := p.Arguments
	dir,largs := largs[0].GetValue(arguments), largs[1:]
	
	cmd := ""
	filestat,error := os.Stat(dir)
	if error != nil || !filestat.IsDir() {
		cmd = dir
		dir = os.Getenv("PWD")
	} else {
		cmd, largs = largs[0].GetValue(arguments), largs[1:]
	}

	args := []string{}
	
	for _,argument := range(largs) {

		args = append(args,argument.GetValue(arguments))
	}

	command := exec.Command(cmd,args...)
	command.Dir = dir
	return command
}

func (p ExecPromise) Desc(arguments []Constant) string {
	cmd := p.getCommand(arguments)
	return "(exec in_dir(" + cmd.Dir + ") <" + cmd.Path + " [" + strings.Join(cmd.Args[1:],", ") + "] >)"
}

func (p ExecPromise) Eval(arguments []Constant, logger *Logger) bool {
	command := p.getCommand(arguments)
	command.Stdout = logger.Stdout
	command.Stderr = logger.Stderr

	logger.Stdout.Write([]byte(strings.Join(command.Args, " ") + "\n"))
	
	err := command.Run()
	
	if (err != nil) {
		return false
	} else {
		return true
	}
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

func (p PipePromise) Eval(arguments []Constant, logger *Logger) bool {
	commands := []*exec.Cmd{}
	cstrings := []string{}
	
	for _,v := range(p.Execs) {
		cmd :=  v.getCommand(arguments)
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
