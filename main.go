package main

import(
	"flag"
	"fmt"
	"os"
	"log"
	"log/syslog"
)

type LogWriter struct {
	log *log.Logger
}

func (l LogWriter) Write(b []byte) (n int, err error) {
	log.Print(string(b))
	return len(b),nil
}

type Command struct {
	Name string
	Usage string
	ShortHelp string
	LongHelp string
	Run func(args []string, logi, loge *log.Logger)
	Flag flag.FlagSet
}

var commands = []*Command{
	run,
	serve,
}

func main() {
	flag.Usage = usage
	var use_syslog = flag.Bool("syslog", false, "output the logs to syslog")
	flag.Parse()
	args := flag.Args()
	
	if len(args) < 1 {
		usage()
		return
	}

	var logi, loge *log.Logger
	
	if *use_syslog {
		logi,_ = syslog.NewLogger(syslog.LOG_NOTICE, log.LstdFlags)
		loge,_ = syslog.NewLogger(syslog.LOG_ERR, log.LstdFlags)
	} else {
		logi = log.New(os.Stdout, "llconf (info)", log.LstdFlags)
		loge = log.New(os.Stderr, "llconf (err)", log.LstdFlags | log.Lshortfile)
	}

	for _,cmd := range commands {
		if cmd.Name == args[0] && cmd.Run != nil {
			cmd.Flag.Parse(args[1:])
			cmd_args := cmd.Flag.Args()
			cmd.Run(cmd_args, logi, loge)
			os.Exit(0)
		}
	}
		
	fmt.Fprintf(os.Stderr, "Unknown subcommand %q\n", args[0])
}

func usage() {
	fmt.Printf("usage: %s\n\n", os.Args[0])

	for _, cmd := range commands {
		if cmd.Usage != "" {
			fmt.Printf("    %s\n",cmd.Usage)
		}
	}
	fmt.Printf("\n")
}
