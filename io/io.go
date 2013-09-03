package io

import (
	"os"
	"io"
	"bufio"
	"errors"
	"path/filepath"
)


func CopyFiles(from, to string) (err error) {
	sf,err := os.Open(from)
	if err != nil { return err }
	defer sf.Close()

	files, err := sf.Readdir(-1)
	if err != nil { return err }
	
	for _,fi := range( files ) {
		src,dest := filepath.Join(from, fi.Name()), filepath.Join(to, fi.Name())
		if fi.IsDir() {
			os.Mkdir(dest,fi.Mode())
			err = CopyFiles(src,dest)
			if err != nil { return err }
		} else {
			err = copyFile(src,dest)
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

type FolderRuneReader struct {
	files []string
	reader io.RuneReader
}

func (this *FolderRuneReader) ReadRune() (r rune, s int, e error) {
	if this.reader == nil {
		if len(this.files) > 0 {
			this.UpdateReader()
		} else {
			return ' ',-1,nil
		}
		
	}


	r,s,e = this.reader.ReadRune()

	if e == io.EOF {
		if len(this.files) != 0 {
			this.UpdateReader()
			r,s,e = this.reader.ReadRune()
		}
	}
	
	return r,s,e
}


func (this *FolderRuneReader) UpdateReader() error {
	for len(this.files) > 0 {
		filename := this.files[0]
		this.files = this.files[1:]
		file,err := os.Open(filename)
		if err != nil {
			continue
		} else {
			this.reader = bufio.NewReader(file)
			return nil
		}
	}
	return errors.New("list of files is empty")
}


func NewFolderRuneReader(folder string) (FolderRuneReader, error) {
	files := []string{}
	
	visit := func(path string, info os.FileInfo, err error) error {
		if ! info.IsDir() {
			files = append(files,path)
		}
		return nil
	}

	err := filepath.Walk(folder, visit)
	if err != nil {
		return FolderRuneReader{}, err
	} else {
		return FolderRuneReader{files: files}, nil
	}
}
