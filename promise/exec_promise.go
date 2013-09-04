package promise

import (
	"os"
	"os/exec"
	"strings"
	"bytes"
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

func (p ExecPromise) Eval(arguments []Constant) (bool,[]string,[]string) {
	command := p.getCommand(arguments)

	var sout,serr bytes.Buffer
	command.Stdout = &sout
	command.Stderr = &serr
	err := command.Run()
	
	stdout := []string{strings.Join(command.Args, " ") + "\n"}
	if str := sout.String(); len(str) > 0 {
		stdout = append(stdout,str)
	}

	stderr := []string{}
	if str := serr.String(); len(str) > 0 {
		stderr = append(stderr,str)
	}
	
	
	if (err != nil) {
		return false,stdout,stderr
	} else {
		return true,stdout,stderr
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

func (p PipePromise) Eval(arguments []Constant) (bool,[]string,[]string) {
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
			panic(err)
		}
		command.Start()
		commands[i + 1].Stdin = out
	}

	var sout,serr bytes.Buffer
	
	last_cmd := commands[len(commands) - 1]
	last_cmd.Stdout = &sout
	last_cmd.Stderr = &serr
	err := last_cmd.Run()

	for _, command := range(commands[:len(commands) - 1]) {
		command.Wait()
	}
	

	stdout := []string{ strings.Join(cstrings, " | ") + "\n" }
	if str := sout.String(); len(str) > 0 {
		stdout = append(stdout,str)
	}

	stderr := []string{}
	if str := serr.String(); len(str) > 0 {
		stderr = append(stderr,str)
	}
		
	if (err != nil) {
		return false,stdout,stderr
	} else {
		return true,stdout,stderr
	}
}
