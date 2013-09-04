package main

import(
	"flag"
	"fmt"
	"os"
	"log"
)

type Command struct {
	Name string
	Usage string
	ShortHelp string
	LongHelp string
	Flag flag.FlagSet
	Run func(logi, loge *log.Logger)
}

var commands = []*Command{
	run,
	serve,
}

var run = &Command{}


func main() {
	flag.Usage = usage
	args := flag.Args()

	if len(args) < 1 {
		usage()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	logi := log.New(os.Stdout, "llconf (info)", log.LstdFlags)
	loge := log.New(os.Stderr, "llconf (err)", log.LstdFlags)

	for _,cmd := range commands {
		if cmd.Name == args[0] && cmd.Run != nil {
			cmd.Flag.Parse(args[1:])
			cmd.Flag.Parse(args[1:])
			cmd.Run(logi, loge)
		}
	}
		
	fmt.Fprintf(os.Stderr, "Unknown subcommand %q\n", args[0])
}

func help(args []string) {
	for _,cmd := range commands {
		if cmd.Name == args[0] {
			fmt.Print(cmd.LongHelp)
		}
	}
}

func usage() {

}
