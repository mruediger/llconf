package main

import (
	"io"
	"os"
	"log"
	"fmt"
	"bufio"
	
	llconf_io "github.com/mruediger/llconf/io"
	"github.com/mruediger/llconf/parse"
	"github.com/mruediger/llconf/promise"
)

var run = &Command{
	Name: "run",
	Usage: "run   [arguments...] [folder]",
	Run: execRun,
}

var run_cfg struct{
	input string
	promise string
	verbose bool
}

func init() {
	run.Flag.BoolVar(&run_cfg.verbose, "verbose", false, "enable verbose output")
	run.Flag.StringVar(&run_cfg.promise, "promise", "done", "the promise that will be used as root")
}

func execRun(args []string, logi, loge *log.Logger) {
	switch len(args) {
	case 0:
		run_cfg.input = ""
		fmt.Println("no input folder specified, reading from stdin")
	case 1:
		run_cfg.input = args[0]
	default:
		os.Exit(1)
	}
	
	input,err := openInput(run_cfg.input)
	if err != nil {
		loge.Printf("could not open %q: %v\n", run_cfg.input, err)
		return
	}
	
	promises,err := parse.ParsePromises(input)
	if err != nil {
		loge.Printf("error while parsing input: %v\n", err)
		return
	}

	p,promise_present := promises[run_cfg.promise]
	if !promise_present {
		loge.Printf("specified goal (%s) not found in config\n", run_cfg.promise)
	}
	success,sout,serr := p.Eval([]promise.Constant{})
	if success {
		if run_cfg.verbose {
			for _,msg := range(sout) {
				logi.Println(msg)
			}
		}
		logi.Println("evaluation successful\n")
	} else {
		var msgs []string
		if run_cfg.verbose {
			msgs = append(sout,serr...)
		} else {
			msgs = serr
		}
		for _,msg := range(msgs) {
			loge.Println(msg)
		}

		loge.Println("error during evaluation")		
	}
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
