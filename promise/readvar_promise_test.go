package promise

import (
	"bytes"
	"testing"
)

func TestReadvarPromise(t *testing.T) {

	vars := Variables{}

	arguments := []Argument{
		Constant{"/"},
		Constant{"/bin/echo"},
		Constant{"Hello World"},
	}
		
		
	exec := ExecPromise{ExecTest, arguments}
	promise := ReadvarPromise{Constant{"test"}, exec}
	logger := NewStdoutLogger()

	var sout bytes.Buffer
	logger.Stdout = &sout
	logger.Info = &sout
	
	promise.Eval([]Constant{},&logger,&vars)
	
	equals(t, "Hello World\n", vars["test"])
	equals(t, "/bin/echo Hello World\nHello World\n", sout.String())

}
