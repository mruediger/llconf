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
	NoInputFolder
)

func (err CliError) Error() string {
	switch err {
	case NotEnoughArguments:
		return "not enough arguments"
	case UnknownMainOption:
		return "unknown main option"
	case NoInputFolder:
		return "no input folder specified"
	default:
		return "unknown error in commandline arguments"
	}
}

type CliConfig interface {
	Run() error
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


func processArguments(args []string) (CliConfig, error) {
	if len(args) < 2 {
		return nil,NotEnoughArguments
	}

	progName := args[0]
	
	switch(args[1]) {
	case "serve":
		return processServeFlags(progName, args[2:])
	case "run":
		return processRunFlags(progName, args[2:])
	default:
		return nil,UnknownMainOption
	}
}

func processServeFlags(progName string, args []string) (CliConfig, error) {
	flagSet := flag.NewFlagSet(progName, 0)
	goal := flagSet.String("promise", "done", "the promise that should be evaluated")
	verbose := flagSet.Bool("verbose", false, "enable verbose output")
	inputFolder := flagSet.String("input-folder", "", "the input folder")
	flagSet.Parse(args)

	if *inputFolder == "" {
		*inputFolder = os.Getenv("LLCONF_INPUT")
		if *inputFolder == "" {
			return nil, NoInputFolder
		}
	}
	
	
	return ServeConfig{ Goal: *goal, InputFolder: *inputFolder, Verbose: *verbose },nil
}

func processRunFlags(progName string, args []string) (CliConfig, error) {
	flagSet := flag.NewFlagSet(progName, 0)
	goal := flagSet.String("promise", "done", "the promise that should be evaluated")
	verbose := flagSet.Bool("verbose", false, "enable verbose output")
	parseOnly := flagSet.Bool("parse-only", false, "only parse the input")
	flagSet.Parse(args)

	input := bufio.NewReader( os.Stdin )
	
	return RunConfig{ Goal: *goal, Input: input, Verbose: *verbose, ParseOnly: *parseOnly },nil
}

func main() {
	runSetup,err := processArguments(os.Args)

	if err == nil {
		runSetup.Run()
	} else {
		fmt.Fprintf(os.Stderr, "argument error: %s\n", err.Error())
	}
}
