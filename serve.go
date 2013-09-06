package main

import (
	"time"
	"log"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/mruediger/llconf/io"
	"github.com/mruediger/llconf/parse"
	libpromise "github.com/mruediger/llconf/promise"
)

var serve = &Command{
	Name: "serve",
	Usage: "serve [arguments...]",
	ShortHelp: "serve",
	LongHelp: "bla",
	Run: runServ,
}

var serve_cfg struct{
	promise string
	verbose bool
	interval int
	inp_dir string
	workdir string
}

func init() {
	serve.Flag.IntVar(&serve_cfg.interval, "interval", 300, "the minium time between promise evaluation")
	serve.Flag.BoolVar(&serve_cfg.verbose, "verbose", false, "enable verbose output")
	serve.Flag.StringVar(&serve_cfg.promise, "promise", "done", "the promise that will be used as root")
	serve.Flag.StringVar(&serve_cfg.inp_dir, "input-folder", "", "the folder containing input files")
}

func runServ(args []string, logi, loge *log.Logger) {
	switch len(args) {
	case 0:
		fmt.Fprintf(os.Stderr, "no workdir specified\n")
		os.Exit(1)
		return
	case 1:
		serve_cfg.workdir = args[0]
	default:
		fmt.Fprintf(os.Stderr, "argument count mismatch")
		os.Exit(1)
	}

	if serve_cfg.inp_dir == "" {
		serve_cfg.inp_dir = filepath.Join(serve_cfg.workdir, "input")
	}
		
	quit := make(chan int)

	var promise libpromise.Promise
	
	for {
		go func(q chan int) {
			time.Sleep(time.Duration(serve_cfg.interval) * time.Second)
			q <- 0
		}(quit)

		new_promise, err := updatePromise(serve_cfg.inp_dir, serve_cfg.promise)
		if err == nil {
			promise = new_promise
		} else {
			loge.Printf("error while parsing input folder: %v\n",err)
		}
		
		checkPromise(promise,logi,loge)
			
		<-quit
	}
}

func updatePromise(folder, root string ) (libpromise.Promise, error) {
	globals := map[string]string{}
	globals["input_dir"] = serve_cfg.inp_dir
	globals["work_dir"] = serve_cfg.workdir
	
	reader,err := io.NewFolderRuneReader( folder )
	if err != nil { return nil, err}

	promises, err := parse.ParsePromises( &reader, &globals )
	if err != nil { return nil, err}

	if promise, promise_present := promises[root]; promise_present {
		return promise, nil
	} else {
		return nil, errors.New("root promise (" + root + ") unknown")
	}
}

func checkPromise(p libpromise.Promise, logi, loge *log.Logger) {
	promises_fullfilled,stdout,stderr := p.Eval([]libpromise.Constant{})

	if promises_fullfilled {
		if serve_cfg.verbose {
			for _,v := range(stdout) {
				logi.Print(v)
			}
		}
		logi.Printf("evaluation successful\n")
	} else {
		var msgs []string
		if serve_cfg.verbose {
			msgs = append(stdout, stderr...)
		} else {
			msgs = stderr
		}
		for _,msg := range(msgs) {
			loge.Print(msg)
		}
		loge.Printf("error during evaluation\n")
	}
}

