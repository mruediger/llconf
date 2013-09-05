package main

import (
	"time"
	"log"

	"github.com/mruediger/llconf/io"
	"github.com/mruediger/llconf/parse"
	"github.com/mruediger/llconf/promise"
)

var serve = &Command{
	Name: "serve",
	Usage: "serve [--verbose]",
	ShortHelp: "serve",
	LongHelp: "bla",
	Run: runServ,
}

var serve_cfg struct{
	promise string
	verbose bool
	interval int
	inc_dir string
	inp_dir string
}

func init() {
	serve.Flag.IntVar(&serve_cfg.interval, "interval", 300, "the minium time between promise evaluation")
	serve.Flag.BoolVar(&serve_cfg.verbose, "verbose", false, "enable verbose output")
	serve.Flag.StringVar(&serve_cfg.promise, "promise", "done", "the promise that should be used as root")
	serve.Flag.StringVar(&serve_cfg.inc_dir, "incomming", "", "the folder containing updates files")
	serve.Flag.StringVar(&serve_cfg.inp_dir, "input", "", "the folder containing input files")
}

func runServ(logi, loge *log.Logger) {
	quit := make(chan int)

	for {
		go func(q chan int) {
			time.Sleep(time.Duration(serve_cfg.interval) * time.Second)
			q <- 0
		}(quit)

		promises,err := parseFolder(serve_cfg.inc_dir)
		if err == nil {
			io.CopyFiles(serve_cfg.inc_dir, serve_cfg.inp_dir)
		} else {
			loge.Printf("error while parsing input folder: %v\n",err)
			promises,err = parseFolder(serve_cfg.inp_dir)
		}

		if promise,promise_present := promises[serve_cfg.promise]; promise_present {
			checkPromises(promise,logi,loge)
		} else {
			loge.Printf("root promise unknown (%s)\n", serve_cfg.promise)
		}
			
		<-quit
	}
}

func parseFolder(folder string) (map[string]promise.Promise, error) {
	reader,err := io.NewFolderRuneReader( folder )
	if err == nil {
		return parse.ParsePromises( &reader )
	} else {
	 	return nil,err
	}
}

func checkPromises(p promise.Promise, logi, loge *log.Logger) {
	promises_fullfilled,stdout,stderr := p.Eval([]promise.Constant{})

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
