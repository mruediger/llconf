package promise

import (
	"bytes"
	"testing"
)

func TestReadvarPromise(t *testing.T) {
	arguments := []Argument{
		Constant{"/bin/echo"},
		Constant{"Hello World"},
	}

	exec := ExecPromise{ExecTest, arguments}
	promise := ReadvarPromise{Constant{"test"}, exec}

	var sout bytes.Buffer

	ctx := NewContext()
	ctx.Logger.Stdout = &sout
	ctx.Logger.Info = &sout

	promise.Eval([]Constant{},&ctx)

	equals(t, "Hello World", ctx.Vars["test"])
	equals(t, "/bin/echo Hello World\nHello World\n", sout.String())

}
