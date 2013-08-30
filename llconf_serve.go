package main

import (
	"os"
	"io"
	"bufio"
	"path/filepath"
	
	"github.com/mruediger/llconf/parse"
	"github.com/mruediger/llconf/promise"
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
	IncommingFolder string
	Verbose bool
}

func (this ServeConfig) Run() error {
	for {
		promises,err := processFolder(this.IncommingFolder)
		if err == nil {
			copyFiles(this.IncommingFolder, this.InputFolder)
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
	reader,err := NewFolderRuneReader( folder )
	if err == nil {
		return parse.ParsePromises( &reader )
	} else {
	 	return nil,err
	}
}

func copyFiles(from, to string) (err error) {
	sf,err := os.Open(from)
	if err != nil { return err }
	defer sf.Close()

	files, err := sf.Readdir(-1)
	if err != nil { return err }
	
	for _,fi := range( files ) {
		if fi.IsDir() {
			continue
		} else {
			err = copyFile( filepath.Join(from, fi.Name()),
				filepath.Join(to, fi.Name()))
			if err != nil { return err }
		}
	}
		
	return nil
}

func copyFile(src_name, dst_name string) (err error) {
	src, err := os.Open(src_name)
	if err != nil { return }
	defer src.Close()

	dst, err := os.Create(dst_name)
	if err != nil { return }
	defer dst.Close()

	_,err = io.Copy(dst,src)
	return
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
