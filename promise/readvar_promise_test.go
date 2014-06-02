package promise

import (
	"bytes"
	"testing"
)

func TestReadvarPromise(t *testing.T) {
	arguments := []Argument{
		Constant("/bin/echo"),
		Constant("Hello World"),
	}

	exec := ExecPromise{ExecTest, arguments}
	promise := ReadvarPromise{Constant("test"), exec}

	var sout bytes.Buffer

	ctx := NewContext()
	ctx.ExecOutput = &sout

	promise.Eval([]Constant{}, &ctx, "readvar")

	equals(t, "Hello World", ctx.Vars["test"])
	equals(t, "Hello World\n", sout.String())
}
