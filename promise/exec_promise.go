package promise

import (
	"os"
	"os/exec"
	"strings"
	"fmt"
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
	return "(exec in_dir(" + cmd.Dir + ") <" + cmd.Path + " [" + strings.Join(cmd.Args,", ") + "] >)"
}

func (p ExecPromise) Eval(arguments []Constant) bool {
	command := p.getCommand(arguments)

	output, e := command.CombinedOutput()

	if e != nil {
		fmt.Println(e, string(output))
		return false
	} else {
		fmt.Println(string(output))
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

func (p PipePromise) Eval(arguments []Constant) bool {
	commands := []*exec.Cmd{}

	for _,v := range(p.Execs) {
		commands = append(commands, v.getCommand(arguments))
	}

	for i, command := range(commands[:len(commands) - 1]) {
		out, err := command.StdoutPipe()
		if err != nil {
			panic(err)
		}
		command.Start()
		commands[i + 1].Stdin = out
	}

	output, err:= commands[len(commands) - 1].Output()
	if (err != nil) {
		fmt.Printf("%v %v\n", err, string(output))
		return false
	} else {
		fmt.Printf("%v\n", string(output))
		return true
	}
}
