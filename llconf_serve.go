package main

import (
	"io"
	"github.com/mruediger/llconf/parse"
)

type ServeConfig struct {
	Goal string
	Input io.RuneReader
	Verbose bool
}

func (this ServeConfig) Run() error {
	for {
	this.processInputFiles()
	}
	
	return nil
}

func (this ServeConfig) processInputFiles() {
	
	parse.ParsePromises( this.Input )
	
}

type FolderRuneReader struct {
	path string
}

func (this FolderRuneReader) ReadRune() (r rune, size int, err error) {
	return ' ',0,nil
}
