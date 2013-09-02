package main

import (
	"io"
	"fmt"

	"github.com/mruediger/llconf/parse"
	"github.com/mruediger/llconf/promise"
)

type RunConfig struct {
	Goal string
	Input io.RuneReader
	Verbose bool
	ParseOnly bool
}

func runClient(cfg RunConfig) error {
	fmt.Println("reading from stdin")
	promises,err := parse.ParsePromises( cfg.Input )

	if err != nil {
		return err
	}

	p,promise_present := promises[cfg.Goal]
	if ! promise_present {
		return SpecifiedGoalUnknown{cfg.Goal}
	}
	
	success,sout,serr := p.Eval([]promise.Constant{})	
	if success {
		fmt.Println("evaluation successful\n")
		if cfg.Verbose {
			for _,v := range(sout) {
				fmt.Print(v)
			}
		}
	} else {
		fmt.Println("error during evaluation\n")
		var msgs []string
		if cfg.Verbose {
			msgs = append(sout,serr...)
		} else {
			msgs = serr
		}
		for _,msg := range(msgs) {
			fmt.Print(msg)
		}
	}

	return nil
}
