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

func (this ServeConfig) Run() error {
	for {
		promises,err := processFolder(this.IncommingFolder)
		if err == nil {
			io.CopyFiles(this.IncommingFolder, this.InputFolder)
		} else {
			log.Printf("error while parsing incomming folder: %v\n", err)
			promises,err = processFolder(this.InputFolder)
		}
		
		success,sout,serr := promises[this.Goal].Eval([]promise.Constant{})	
		if success {
			log.Printf("evaluation successful\n")
			if this.Verbose {
				for _,v := range(sout) {
					log.Print(v)
				}
			}
		} else {
			log.Printf("error during evaluation\n")
			var msgs []string
			if this.Verbose {
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
