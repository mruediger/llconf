package main

import (
	"fmt"
	"os"
	"bufio"
	"flag"
)

type CliError int

const (
	NotEnoughArguments CliError = iota
	UnknownMainOption
	NoInputFolder
	NoIncommingFolder
)

func (err CliError) Error() string {
	switch err {
	case NotEnoughArguments:
		return "not enough arguments"
	case UnknownMainOption:
		return "unknown main option"
	case NoInputFolder:
		return "no input folder was specified"
	case NoIncommingFolder:
		return "no incomming folder for downloaded files was specified"
	default:
		return "unknown error in commandline arguments"
	}
}


type SpecifiedGoalUnknown struct  {
	Goal string
}

func (err SpecifiedGoalUnknown) Error() string {
	return "specified goal (" + err.Goal +") not found in config"
}

func usage() {
	fmt.Printf("usage: %s\n", os.Args[0])
	fmt.Printf("\n")
	fmt.Printf("   serve   starts serving files\n")
	fmt.Printf("   run     evaluates promises from stdin\n")
	fmt.Printf("\n")
}

func main() {
	args := os.Args
	
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "%s\n", NotEnoughArguments.Error())
		usage()
		os.Exit(1)
	}

	progName := args[0]

		
	switch(args[1]) {
	case "serve":
		config,err := processServeFlags(progName, args[2:])
		if err == nil {
			runerr := runServer(config)
			if runerr != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", runerr.Error())
			}
		} else {
			fmt.Fprintf(os.Stderr, "argument error: %s\n", err.Error())
			os.Exit(1)
		}
	case "run":
		config,err := processRunFlags(progName, args[2:])
		if err == nil {
		 	runerr := runClient(config)
			if runerr != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", runerr.Error())
			}
		} else {
			fmt.Fprintf(os.Stderr, "argument error: %s\n", err.Error())
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[1])
		usage()
		os.Exit(1)
	}
}


func processServeFlags(progName string, args []string) (ServeConfig, error) {
	flagSet := flag.NewFlagSet(progName, flag.ExitOnError)
	goal := flagSet.String("promise", "done", "the promise that should be evaluated")
	verbose := flagSet.Bool("verbose", false, "enable verbose output")
	inputFolder := flagSet.String("input-folder", "", "the folder for input files")
	incommingFolder := flagSet.String("incomming-folder", "", "the folder for updates")

	flagSet.Usage = func() {
		fmt.Printf("usage: %s serve\n\n", os.Args[0])
		flagSet.PrintDefaults()
		fmt.Printf("\n")
	}
	
	flagSet.Parse(args)

	if *inputFolder == "" {
		flagSet.Usage()			
		return ServeConfig{}, NoInputFolder
	}

	if *incommingFolder == "" {
		flagSet.Usage()
		return ServeConfig{}, NoIncommingFolder
	}
	return ServeConfig{ Goal: *goal, InputFolder: *inputFolder, Verbose: *verbose },nil
}

func processRunFlags(progName string, args []string) (RunConfig, error) {
	flagSet := flag.NewFlagSet(progName, flag.ExitOnError)
	goal := flagSet.String("promise", "done", "the promise that should be evaluated")
	verbose := flagSet.Bool("verbose", false, "enable verbose output")
	parseOnly := flagSet.Bool("parse-only", false, "only parse the input")

	flagSet.Usage = func() {
		fmt.Printf("usage: %s run\n\n", os.Args[0])
		flagSet.PrintDefaults()
		fmt.Printf("\n")
	}
	
	flagSet.Parse(args)

	input := bufio.NewReader( os.Stdin )
	
	return RunConfig{ Goal: *goal, Input: input, Verbose: *verbose, ParseOnly: *parseOnly },nil
}

