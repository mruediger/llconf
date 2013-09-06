package promise

import (
	"testing"
)

func TestAndPromiseDesc(t *testing.T) {
	promise := AndPromise{ []Promise{ DummyPromise{ "hello world", true } } }
	desc := promise.Desc([]Constant{})
	equals(t, desc,"(and (dummy [hello world]))")
}

func TestAndPromiseEvalAllTrue(t *testing.T) {
	promise := AndPromise{ []Promise{} }

	for i :=0; i < 10; i++ {
		promise.Promises = append(promise.Promises, DummyPromise{ "n:" + string(1), true })
	}

	result := promise.Eval([]Constant{}, &Logger{})
	equals(t, true, result)
}

func TestAndPromiseEvalSomeFalse(t *testing.T) {
	promise := AndPromise{ []Promise{} }
	promise.Promises = append(promise.Promises, DummyPromise{ "n:0", false })
	for i :=1; i < 10; i++ {
		promise.Promises = append(promise.Promises, DummyPromise{ "n:" + string(1), true })
	}

	result := promise.Eval([]Constant{}, &Logger{})
	equals(t, false, result)
}
