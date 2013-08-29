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

type CliError int

const (
	NotEnoughArguments CliError = iota
	UnknownMainOption
)

func (err CliError) Error() string {
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
	Input io.RuneReader
	Verbose bool
}

func (this ServeConfig) Run() error {
	return nil
}

type RunConfig struct {
	Goal string
	Input io.RuneReader
	Verbose bool
	ParseOnly bool
}

func (this RunConfig) Run() error {
	fmt.Println("reading from stdin")
	promises,err := parse.ParsePromises( this.Input )

	success := promises[this.Goal].Eval([]promise.Constant{})

	if success {
		fmt.Println("evaluation successful")
	} else {
		fmt.Println("evaluation not successful")
	}

	return err
}


func processArguments(args []string, input io.RuneReader) (CliConfig, error) {
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

func processServeFlags(progName string, input io.RuneReader,  args []string) (CliConfig, error) {
	flagSet := flag.NewFlagSet(progName, 0)
	goal := flagSet.String("promise", "done", "the promise that should be evaluated")
	verbose := flagSet.Bool("verbose", false, "enable verbose output")
	flagSet.Parse(args)

	return ServeConfig{ Goal: *goal, Input: input, Verbose: *verbose, ParseOnly: *parseOnly },nil
}

func processRunFlags(progName string, input io.RuneReader, args []string) (CliConfig, error) {
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
