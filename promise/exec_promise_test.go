package promise

import (
	"bytes"
	"strconv"
	"testing"
)

func TestExecPromise(t *testing.T) {
	var promise Promise
	promise = ExecPromise{ExecTest, []Argument{
		Constant("/bin/echo"),
		Constant("Hello"),
		ArgGetter{0}}}


	out := bytes.Buffer{}
	ctx := NewContext()
	ctx.ExecOutput = &out

	res := promise.Eval([]Constant{}, &ctx)
	equals(t, true, res)

	equals(t, strconv.Itoa(7), strconv.Itoa(len(out.String())))
	equals(t, "Hello \n", out.String())
}

func TestPipePromise(t *testing.T) {
	exec1 := ExecPromise{ExecTest, []Argument{Constant("/bin/echo"),
		Constant("hello world")}}
	exec2 := ExecPromise{ExecChange, []Argument{Constant("/usr/bin/rev")}}

	var promise Promise

	promises := []ExecPromise{exec1, exec2}
	promise = PipePromise{promises}

	var out bytes.Buffer
	ctx := NewContext()
	ctx.ExecOutput = &out;

	res := promise.Eval([]Constant{}, &ctx)
	equals(t, true, res)
	equals(t, strconv.Itoa(12), strconv.Itoa(len(out.String())))
	equals(t, "dlrow olleh\n", out.String())
}

func TestExecReporting(t *testing.T) {
	arguments := []Argument{
		Constant("/bin/echo"),
		Constant("Hello"),
		ArgGetter{0},
	}

	var tests = []struct {
		promise ExecPromise
		changes int
	}{
		{ExecPromise{ExecChange, arguments}, 1},
		{ExecPromise{ExecTest, arguments}, 0},
	}

	for _, test := range tests {
		var out bytes.Buffer
		ctx := NewContext()
		ctx.ExecOutput = &out;

		res := test.promise.Eval([]Constant{}, &ctx)
		equals(t, true, res)
		equals(t, strconv.Itoa(test.changes), strconv.Itoa(ctx.Logger.Changes))
	}
}
