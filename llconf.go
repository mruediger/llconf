package main

import (
	"fmt"
	"os"
	"bufio"
	"flag"
	"io"
	"github.com/mruediger/llconf/parse"
	"github.com/mruediger/llconf/promise"
)

type ArgumentError int

const (
	NotEnoughArguments ArgumentError = iota
	UnknownMainOption
)

func (err ArgumentError) Error() string {
	switch err {
	case NotEnoughArguments:
		return "not enough arguments"
	case UnknownMainOption:
		return "unknown main option"
	default:
		return "unknown error in commandline arguments"
	}
}

type CliConfig interface {
	Run() error
}

type ServeConfig struct {
	Goal string
	Input io.RuneScanner
	Verbose bool
	ParseOnly bool
}

func (this ServeConfig) Run() error {
	return nil
}

type RunConfig struct {
	Goal string
	Input io.RuneScanner
	Verbose bool
	ParseOnly bool
}

func (this RunConfig) Run() error {
	fmt.Println("reading from stdin")
	promises := parse.ParsePromises( this.Input )
	for k,v := range(promises) {
		fmt.Printf("present %s\n", v)

		if k == this.Goal {
			fmt.Printf("evaluating %s\n", v)
			if v.Eval([]promise.Constant{}) {
				fmt.Println("success")
			} else {
				fmt.Println("failure")
			}
		} 
	}
	return nil
}


func processArguments(args []string, input io.RuneScanner) (CliConfig, error) {
	if len(args) < 2 {
		return nil,NotEnoughArguments
	}

	progName := args[0]
	
	switch(args[1]) {
	case "serve":
		return processServeFlags(progName, input, args[2:])
	case "run":
		return processRunFlags(progName, input, args[2:])
	default:
		return nil,UnknownMainOption
	}
}

func processServeFlags(progName string, input io.RuneScanner,  args []string) (CliConfig, error) {
	flagSet := flag.NewFlagSet(progName, 0)
	goal := flagSet.String("promise", "done", "the promise that should be evaluated")
	verbose := flagSet.Bool("verbose", false, "enable verbose output")
	parseOnly := flagSet.Bool("parse-only", false, "only parse the input")
	flagSet.Parse(args)

	return ServeConfig{ Goal: *goal, Input: input, Verbose: *verbose, ParseOnly: *parseOnly },nil
}

func processRunFlags(progName string, input io.RuneScanner, args []string) (CliConfig, error) {
	flagSet := flag.NewFlagSet(progName, 0)
	goal := flagSet.String("promise", "done", "the promise that should be evaluated")
	verbose := flagSet.Bool("verbose", false, "enable verbose output")
	parseOnly := flagSet.Bool("parse-only", false, "only parse the input")
	flagSet.Parse(args)

	return RunConfig{ Goal: *goal, Input: input, Verbose: *verbose, ParseOnly: *parseOnly },nil
}

func main() {
	in := os.Stdin
	bufin := bufio.NewReader( in )

	runSetup,err := processArguments(os.Args, bufin)

	if err == nil {
		runSetup.Run()
	} else {
		fmt.Fprintf(os.Stderr, "argument error: %s\n", err.Error())
//		usage()
	}
}
