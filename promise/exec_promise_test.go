package promise

import (
	"testing"
)

func TestExecPromiseDesc(t *testing.T) {
	var promise Promise
	promise = ExecPromise{[]Argument{Constant{"/"},
		Constant{"/bin/echo"},
		Constant{"Hello"},
		ArgGetter{0}}}
	desc := promise.Desc([]Constant{Constant{"World"}})
	equals(t, "(exec in_dir(/) </bin/echo [Hello, World] >)", desc)

	res,sout,serr := promise.Eval([]Constant{})
	equals(t, true, res)
	equals(t, string(2), string(len(sout)))
	equals(t, "Hello \n", sout[1])
	equals(t, string(0), string(len(serr)))
}

func TestPipePromiseDesc(t *testing.T) {
	exec1 := ExecPromise{[]Argument{Constant{"/tmp"},
		Constant{"/bin/echo"},
		Constant{"hello world"}}}
	exec2 := ExecPromise{[]Argument{Constant{"/tmp"},
		Constant{"/usr/bin/rev"}}}

	var promise Promise
	
	promises := []ExecPromise{exec1, exec2}
	promise = PipePromise{promises}
	desc := promise.Desc([]Constant{})
	equals(t, "(pipe "+
		"(exec in_dir(/tmp) </bin/echo [hello world] >) "+
		"(exec in_dir(/tmp) </usr/bin/rev [] >))", desc)

	res,sout,serr := promise.Eval([]Constant{})
	equals(t, true, res)
	equals(t, string(2), string(len(sout)))
	equals(t, "dlrow olleh\n", sout[1])
	equals(t, string(0), string(len(serr)))
}
