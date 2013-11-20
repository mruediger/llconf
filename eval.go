package main

import (
	"io"
	"os"
	"fmt"
	"bufio"

	llconf_io "github.com/mruediger/llconf/io"
	"github.com/mruediger/llconf/parse"

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

	input,err := openInput(run_cfg.input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open %q: %v\n", run_cfg.input, err)
		return
	}

	promises,err := parse.ParsePromises(input)
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

func openInput( source string ) (io.RuneReader, error) {
	if source == "" {
		input := bufio.NewReader( os.Stdin )
		return input,nil
	} else {
		input,err := llconf_io.NewFolderRuneReader( run_cfg.input )
		return &input,err
	}
}
