package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"path/filepath"
	"time"

	"bitbucket.org/kardianos/osext"

	"github.com/d3media/llconf/compiler"
	libpromise "github.com/d3media/llconf/promise"
	"bytes"
)

var serve = &Command{
	Name:      "serve",
	Usage:     "serve [arguments...]",
	ShortHelp: "serve",
	LongHelp:  "bla",
	Run:       runServ,
}

var serve_cfg struct {
	root_promise string
	verbose      bool
	use_syslog   bool
	interval     int
	inp_dir      string
	workdir      string
	runlog_path  string
}

func init() {
	serve.Flag.IntVar(&serve_cfg.interval, "interval", 300, "set the minium time between promise-tree evaluation")
	serve.Flag.BoolVar(&serve_cfg.verbose, "verbose", false, "enable verbose output")
	serve.Flag.StringVar(&serve_cfg.root_promise, "promise", "done", "the promise that will be used as the root of the promise tree")
	serve.Flag.StringVar(&serve_cfg.inp_dir, "input-folder", "", "the folder containing input files")
	serve.Flag.BoolVar(&serve_cfg.use_syslog, "syslog", false, "output to syslog")
	serve.Flag.StringVar(&serve_cfg.runlog_path, "runlog", "", "path to the runlog")
}

func runServ(args []string) {
	parseArguments(args)
	logi, loge := setupLogging()

	quit := make(chan int)

	var promise_tree libpromise.Promise

	for {
		go func(q chan int) {
			time.Sleep(time.Duration(serve_cfg.interval) * time.Second)
			q <- 0
		}(quit)

		new_promise_tree, err := updatePromise(serve_cfg.inp_dir, serve_cfg.root_promise)
		if err == nil {
			promise_tree = new_promise_tree
		} else {
			loge.Printf("error while parsing input folder: %v\n", err)
		}

		if promise_tree != nil {
			checkPromise(promise_tree, logi, loge, flag.Args())
		} else {
			fmt.Fprintf(os.Stderr, "could not find any valid promises\n")
		}
		<-quit
	}
}

func setupLogging() (logi, loge *log.Logger) {
	if serve_cfg.use_syslog {
		logi, _ = syslog.NewLogger(syslog.LOG_NOTICE, log.LstdFlags)
		loge, _ = syslog.NewLogger(syslog.LOG_ERR, log.LstdFlags)
		return
	} else {
		logi = log.New(os.Stdout, "llconf (info)", log.LstdFlags)
		loge = log.New(os.Stderr, "llconf (err)", log.LstdFlags|log.Lshortfile)
		return
	}
}

func parseArguments(args []string) {
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

	if serve_cfg.runlog_path == "" {
		serve_cfg.runlog_path = filepath.Join(serve_cfg.workdir, "runlog")
	}

	// when run as daemon, the home folder isn't set
	home := os.Getenv("HOME")
	if home == "" {
		os.Setenv("HOME", serve_cfg.workdir)
	}
}

func updatePromise(folder, root string) (libpromise.Promise, error) {
	promises, err := compiler.Compile(folder)
	if err != nil {
		return nil, err
	}

	if promise, promise_present := promises[root]; promise_present {
		return promise, nil
	} else {
		return nil, errors.New("root promise (" + root + ") unknown")
	}
}

func checkPromise(p libpromise.Promise, logi, loge *log.Logger, args []string) {
	vars := libpromise.Variables{}
	vars["input_dir"] = serve_cfg.inp_dir
	vars["work_dir"] = serve_cfg.workdir
	exe, _ := osext.Executable()
	vars["executable"] = exe
	env := []string{}

	logger := libpromise.Logger{
		Error:   loge,
		Info:    logi,
		Changes: 0,
		Tests:   0}

	ctx := libpromise.Context{
		Logger: &logger,
		ExecOutput: &bytes.Buffer{},
		Vars:   vars,
		Args:   args,
		Env:    env,
		Debug:  false,
		InDir:  ""}

	starttime := time.Now().Local()
	promises_fullfilled := p.Eval([]libpromise.Constant{}, &ctx)
	endtime := time.Now().Local()

	logi.Printf("%d changes and %d tests executed\n", ctx.Logger.Changes, ctx.Logger.Tests)
	if promises_fullfilled {
		logi.Printf("evaluation successful\n")
	} else {
		loge.Printf("error during evaluation\n")
	}

	writeRunLog(promises_fullfilled, ctx.Logger.Changes, ctx.Logger.Tests, starttime, endtime, serve_cfg.runlog_path)
}

func writeRunLog(success bool, changes, tests int, starttime, endtime time.Time, path string) (err error) {
	var output string

	duration := endtime.Sub(starttime)

	if success {
		output = fmt.Sprintf("successful, endtime=%d, duration=%f, c=%d, t=%d\n", endtime.Unix(), duration.Seconds(), changes, tests)
	} else {
		output = fmt.Sprintf("error, endtime=%d, duration=%f, c=%d, t=%d\n", endtime.Unix(), duration.Seconds(), changes, tests)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return
	}

	data := []byte(output)
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		return
	}

	err = f.Close()
	return
}
