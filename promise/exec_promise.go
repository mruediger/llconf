package promise

import (
	"os"
	"os/exec"
	"strings"
	"fmt"
	"bytes"
	"io"
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
		fmt.Println(output)
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
	cmds := []*exec.Cmd{}

	for _,v := range(p.Execs) {
		cmds = append(cmds, v.getCommand(arguments))
	}

	prev_cmd := cmds[0]
	for _,cmd := range(cmds[1:]) {
		pipe, _ := cmd.StdinPipe()
		prev_cmd.Stdout = pipe
	}

	last_cmd,cmds := cmds[len(cmds)-1], cmds[:len(cmds)-1]
	var output bytes.Buffer
	
	last_cmd.Stdout = &output
	last_cmd.Start()

	for _,cmd := range(cmds) {
		err := cmd.Start()
		if err != nil {
			panic(err)
		}
	}

	for _,cmd := range(cmds) {
		err := cmd.Wait()
		cmd.Stdout.(io.WriteCloser).Close()
		if err != nil {
			panic(err)
		}
	}

	err := last_cmd.Wait()
	if err != nil {
		fmt.Printf("%v %v\n", err, output.String())
		return false
	} else {
		fmt.Printf("%v\n",output.String())
		return true
	}
	return false
}
