package promise

import "testing"

func TestNamedPromiseDesc(t *testing.T) {
	promise := NamedPromise{ "test", DummyPromise{ "Hello", true } }
	equals(t, promise.Desc([]Constant{}), "(test (dummy [Hello]))")
}

func TestNamedPromiseEval(t *testing.T) {
	promise := NamedPromise{ "test", DummyPromise{ "Hello", true }}
	result := promise.Eval([]Constant{}, &Logger{})
	equals(t, result, true)
}
