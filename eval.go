package main

import (
	"os"
	"fmt"

	"github.com/mruediger/llconf/compiler"

)

var eval = &Command{
	Name: "eval",
	Usage: "eval [input_folder]",
	Run: evalRun,
}

var run_cfg struct{
	input string
	promise string
	verbose bool
	dryrun bool
}

func init() {
	eval.Flag.StringVar(&run_cfg.promise, "promise", "done", "the promise that will be used as root")
}

func evalRun(args []string) {
	switch len(args) {
	case 0:
		fmt.Println("no input folder specified")
		os.Exit(1)
	case 1:
		run_cfg.input = args[0]
	default:
		fmt.Fprintf(os.Stderr, "argument count mismatch")
		os.Exit(1)
	}

	promises,err := compiler.Compile(run_cfg.input)
	fmt.Println(promises)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while parsing input: %v\n", err)
		return
	}

	_,promise_present := promises[run_cfg.promise]
	if !promise_present {
		fmt.Fprintf(os.Stderr, "specified goal (%s) not found in config\n", run_cfg.promise)
		return
	}
	fmt.Println("evaluation successfull")
}
