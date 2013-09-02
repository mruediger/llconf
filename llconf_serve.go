package main

import (
	"log"
	
	"github.com/mruediger/llconf/io"
	"github.com/mruediger/llconf/parse"
	"github.com/mruediger/llconf/promise"
)

type ServeConfig struct {
	Goal string
	InputFolder string
	IncommingFolder string
	Verbose bool
}

func runServer(cfg ServeConfig) error {
	for {
		promises,err := processFolder(cfg.IncommingFolder)
		if err == nil {
			io.CopyFiles(cfg.IncommingFolder, cfg.InputFolder)
		} else {
			log.Printf("error while parsing incomming folder: %v\n", err)
			promises,err = processFolder(cfg.InputFolder)
		}

		p, promise_present := promises[cfg.Goal]
		if ! promise_present {
			return SpecifiedGoalUnknown{cfg.Goal}
		}
		
		success,sout,serr := p.Eval([]promise.Constant{})	
		if success {
			log.Printf("evaluation successful\n")
			if cfg.Verbose {
				for _,v := range(sout) {
					log.Print(v)
				}
			}
		} else {
			log.Printf("error during evaluation\n")
			var msgs []string
			if cfg.Verbose {
				msgs = append(sout,serr...)
			} else {
				msgs = serr
			}
			for _,msg := range(msgs) {
				log.Print(msg)
			}
		}
	}
	
	return nil
}

func processFolder(folder string) (map[string]promise.Promise, error) {
	reader,err := io.NewFolderRuneReader( folder )
	if err == nil {
		return parse.ParsePromises( &reader )
	} else {
	 	return nil,err
	}
}
