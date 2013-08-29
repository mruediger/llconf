package main

import (
	"os"
	"io"
	"bufio"

	"path/filepath"
	
	"github.com/mruediger/llconf/parse"
)

type FolderRuneReaderError int

const (
	InputFolderIsNoDir FolderRuneReaderError = iota
)

func (this FolderRuneReaderError) Error() string {
	switch this {
	case InputFolderIsNoDir:
		return "the input folder is not a folder"
	default:
		return "unknown error"
	}
}

type ServeConfig struct {
	Goal string
	InputFolder string
	Verbose bool
}

func (this ServeConfig) Run() error {
	for {
		this.processInputFiles()
	}
	
	return nil
}

func (this ServeConfig) processInputFiles() {
	reader,err := NewFolderRuneReader( this.InputFolder )
	if err == nil {
		parse.ParsePromises( &reader )
	} else {
		panic(err)
	}
	
}

type FolderRuneReader struct {
	files []string
	reader io.RuneReader
}

func (this *FolderRuneReader) ReadRune() (r rune, s int, e error) {
	
	r,s,e = this.reader.ReadRune()

	if e == io.EOF {
		if len(this.files) != 0 {
			filename := this.files[0]
			this.files = this.files[1:]
			file,_ := os.Open(filename)
			this.reader = bufio.NewReader( file )
			r,s,e = this.reader.ReadRune()
		}
	}
	return r,s,e
}

func NewFolderRuneReader(folder string) (FolderRuneReader, error) {
	files := []string{}
	fp,err := os.Open(folder)
	basenames,err := fp.Readdir(-1)
	if err != nil {
		return FolderRuneReader{},err
	}

	for _,v := range( basenames ) {
		if ! v.IsDir() {
			files = append( files, filepath.Join(fp.Name(), v.Name()))
		}
	}

	fp,err = os.Open(files[0])
	
	if err != nil {
		return FolderRuneReader{},err
	}
	
	reader := bufio.NewReader( fp )		
	
	files = files[1:]
	return FolderRuneReader{files, reader }, nil
	
}
