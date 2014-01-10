package compiler

import (
	"fmt"
	"testing"
)

func TestListFiles(t *testing.T) {
	ch := make(chan string)

	go listFiles("/home/mathias/d3media/llconf/", "cnf", ch)

	found := 0
	for file := range ch {
		fmt.Println(file)
		found ++
	}

	if found == 0 {
		t.Errorf("didn't found any file")
	}

}
