package main

import (
	"time"
	"log"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"log/syslog"
	
	"github.com/mruediger/llconf/io"
	"github.com/mruediger/llconf/parse"
	libpromise "github.com/mruediger/llconf/promise"
)

type LogWriter struct {
	log *log.Logger
}

func (l LogWriter) Write(b []byte) (n int, err error) {
	l.log.Print(string(b))
	return len(b),nil
}

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
	use_syslog bool
	interval int
	inp_dir string
	workdir string
}

func init() {
	serve.Flag.IntVar(&serve_cfg.interval, "interval", 300, "the minium time between promise evaluation")
	serve.Flag.BoolVar(&serve_cfg.verbose, "verbose", false, "enable verbose output")
	serve.Flag.StringVar(&serve_cfg.promise, "promise", "done", "the promise that will be used as root")
	serve.Flag.StringVar(&serve_cfg.inp_dir, "input-folder", "", "the folder containing input files")
	serve.Flag.BoolVar(&serve_cfg.use_syslog, "syslog", false, "log to syslog instead of to stdout")
}

func runServ(args []string) {
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

	// when run as daemon, the home folder isn't set
    home := os.Getenv("HOME")
	if home == "" {
		os.Setenv("HOME", serve_cfg.workdir)
	}
		
	var logi, loge *log.Logger

	if serve_cfg.use_syslog {
		fmt.Println("using syslog")
		logi,_ = syslog.NewLogger(syslog.LOG_NOTICE, log.LstdFlags)
		loge,_ = syslog.NewLogger(syslog.LOG_ERR, log.LstdFlags)
	} else {
		logi = log.New(os.Stdout, "llconf (info)", log.LstdFlags)
		loge = log.New(os.Stderr, "llconf (err)", log.LstdFlags | log.Lshortfile)
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
	logger := libpromise.Logger{ LogWriter{ logi }, LogWriter{ loge } }
	promises_fullfilled := p.Eval([]libpromise.Constant{}, &logger)

	if promises_fullfilled {
		logi.Printf("evaluation successful\n")
	} else {
		loge.Printf("error during evaluation\n")
	}
}

