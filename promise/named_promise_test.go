package promise

import "testing"

func TestNamedPromiseDesc(t *testing.T) {
	promise := NamedPromise{"test", DummyPromise{"Hello", true}, []Argument{}}
	equals(t, promise.Desc([]Constant{}), "(test (dummy [Hello]))")
}

func TestNamedPromiseEval(t *testing.T) {
	promise := NamedPromise{"test", DummyPromise{"Hello", true}, []Argument{}}
	result := promise.Eval([]Constant{}, &Context{}, "namedpromise")
	equals(t, result, true)
}
