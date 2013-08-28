package promise

import (
	"testing"
)

func TestExecPromiseDesc(t *testing.T) {
	promise := ExecPromise{[]Argument{Constant{"/"},
		Constant{"/bin/echo"},
		Constant{"Hello"},
		ArgGetter{0}}}
	desc := promise.Desc([]Constant{Constant{"World"}})
	equals(t, "(exec in_dir(/) </bin/echo [Hello, World] >)", desc)
}

func TestPipePromiseDesc(t *testing.T) {
	exec1 := ExecPromise{[]Argument{Constant{"/tmp"},
		Constant{"/bin/echo"},
		Constant{"hello world"}}}
	exec2 := ExecPromise{[]Argument{Constant{"/tmp"},
		Constant{"/usr/bin/rev"}}}

	promises := []ExecPromise{exec1, exec2}
	promise := PipePromise{promises}
	desc := promise.Desc([]Constant{})
	equals(t, "(pipe "+
		"(exec in_dir(/tmp) </bin/echo [hello world] >) "+
		"(exec in_dir(/tmp) </usr/bin/rev [] >))", desc)
}
