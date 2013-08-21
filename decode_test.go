package llconf

import (
	"testing"
	"os"
	"bytes"
	"fmt"
)

func TestOpenFile(t *testing.T) {
	content, error := openFile("./setup.cp")

	if error != nil {
		fmt.Println(error)
	} else {
		fmt.Println(content)
	}
}

func openFile(filename string) (string, error) {
	file, err := os.Open(filename)
	buf := new(bytes.Buffer)


	buf.ReadFrom(file)


	return buf.String(), err
}
