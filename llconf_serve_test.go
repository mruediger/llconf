package main

import (
	"testing"
	"io"
	"io/ioutil"
	"os"
)


func TestNewFolderRuneReader(t *testing.T) {
	reader,err := NewFolderRuneReader(".")
	if err != nil {
		t.Errorf("%v\n",err)
	}

	found_file := false
	filename := "llconf_serve_test.go"
	for _,v := range( reader.files ) {
		if v == filename {
			found_file = true
		}
	}

	if !found_file {
		t.Errorf("didn't found " + filename )
	}
	
	for {
		_,_,err := reader.ReadRune()
		if err != nil {
			if err != io.EOF {
				t.Errorf("%v\n",err)
			}
			break
		}
	}

}

func TestCopyFiles(t *testing.T) {

	df,dfe := ioutil.TempDir("/tmp", "copy-from")
	dt,dte := ioutil.TempDir("/tmp", "copy-to")

	if dfe != nil { panic(dfe) }
	if dte!= nil { panic(dte) }

	files := []string{}

	for i := 0; i < 10; i++ {
		fh,e := ioutil.TempFile(df, "copy-test")
		fh.Close()
		files = append(files,fh.Name())
		if e != nil { panic(e) }
	}
	

	
	err := copyFiles(df,dt)
	if err != nil {
		t.Errorf("error while copying: %v\n", err)
	}

	for _,file := range(files) {
		fh,e := os.Open(file)
		if err != nil {
			t.Errorf("error opening file: %v\n", e)
		}
		fh.Close()
	}
	
	os.RemoveAll(df)
	os.RemoveAll(dt)
}
