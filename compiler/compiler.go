package compiler

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mruediger/llconf/compiler/parser"
	"github.com/mruediger/llconf/promise"
)

type FolderReader struct {
}

func Compile(folder string) (map[string]promise.Promise, error) {
	ch := make(chan string)
	go listFiles(folder, "cnf", ch)

	inputs := []parser.Input{}

	for filename := range ch {
		if content, err := ioutil.ReadFile(filename); err != nil {
			return nil, err
		} else {
			inputs = append(inputs, parser.Input{
				filename,
				string(content)})
		}
	}

	return parser.Parse(inputs)
}

func listFiles(folder, suffix string, filename chan<- string) {
	defer close(filename)

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), suffix) {
			filename <- path
		}
		return nil
	})

	if err != nil {
		fmt.Printf(err.Error())
	}
}
