package main

import (
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
			promises,err = processFolder(this.InputFolder)

			// TODO report error
		}

		
		promises[this.Goal].Eval([]promise.Constant{})	
		//TODO report result
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
