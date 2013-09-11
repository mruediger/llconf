package promise

import (
	"testing"
	"bytes"
	"strconv"
)

func TestExecPromiseDesc(t *testing.T) {
	var promise Promise
	promise = ExecPromise{ExecTest,[]Argument{Constant{"/"},
		Constant{"/bin/echo"},
		Constant{"Hello"},
		ArgGetter{0}}}
	desc := promise.Desc([]Constant{Constant{"World"}})
	equals(t, "(test in_dir(/) </bin/echo [Hello, World] >)", desc)

	var sout,serr bytes.Buffer
	logger := Logger{Stdout:&sout, Stderr:&serr}
	
	res := promise.Eval([]Constant{},&logger)
	equals(t, true, res)
	equals(t, strconv.Itoa(24), strconv.Itoa(len(sout.String())))
	equals(t, "/bin/echo Hello \nHello \n", sout.String())
	equals(t, strconv.Itoa(0), strconv.Itoa(len(serr.String())))
}

func TestPipePromiseDesc(t *testing.T) {
	exec1 := ExecPromise{ExecTest, []Argument{Constant{"/tmp"},
		Constant{"/bin/echo"},
		Constant{"hello world"}}}
	exec2 := ExecPromise{ExecChange, []Argument{Constant{"/tmp"},
		Constant{"/usr/bin/rev"}}}

	var promise Promise
	
	promises := []ExecPromise{exec1, exec2}
	promise = PipePromise{promises}
	desc := promise.Desc([]Constant{})
	equals(t, "(pipe "+
		"(test in_dir(/tmp) </bin/echo [hello world] >) "+
		"(change in_dir(/tmp) </usr/bin/rev [] >))", desc)

	var sout,serr bytes.Buffer
	logger := Logger{Stdout:&sout, Stderr:&serr}
		
	res := promise.Eval([]Constant{}, &logger)
	equals(t, true, res)
	equals(t, strconv.Itoa(49), strconv.Itoa(len(sout.String())))
	equals(t, "/bin/echo hello world | /usr/bin/rev\ndlrow olleh\n", sout.String())
	equals(t, strconv.Itoa(0), strconv.Itoa(len(serr.String())))
}

func TestExecReporting(t *testing.T) {
	arguments := []Argument {
		Constant{"/"},
		Constant{"/bin/echo"},
		Constant{"Hello"},
		ArgGetter{0},
	}

	var tests = []struct {
		promise ExecPromise
		changes int
	}{
		{ExecPromise{ExecChange, arguments},1},
		{ExecPromise{ExecTest, arguments},0},
	}
	
	for _,test := range tests {
		var sout,serr bytes.Buffer
		logger := Logger{Stdout:&sout, Stderr:&serr}
		res := test.promise.Eval([]Constant{},&logger)
		desc := test.promise.Desc([]Constant{Constant{"World"}})
		equals(t, true, res)
		equals(t, "(" + test.promise.Type.Name() + " in_dir(/) </bin/echo [Hello, World] >)", desc)
		
		equals(t, strconv.Itoa(test.changes), strconv.Itoa(len(logger.Changes)))
	}
}
