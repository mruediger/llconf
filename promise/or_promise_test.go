package promise

import (
	"testing"
)

func TestOrPromiseDesc(t *testing.T) {
	promise := OrPromise{ []Promise{ DummyPromise{ "hello world", true } }}

	desc := promise.Desc([]Constant{})
	equals(t, desc,"(or (dummy [hello world]))")
}

func TestOrPromiseEvalAllTrue(t *testing.T) {
	promise := OrPromise{ []Promise{} }

	for i :=0; i < 10; i++ {
		promise.Promises = append(promise.Promises, DummyPromise{ "n:" + string(1), true })
	}

	result,_,_ := promise.Eval([]Constant{})
	equals(t, true, result)
}

func TestOrPromiseEvalSomeFalse(t *testing.T) {
	promise := OrPromise{ []Promise{} }
	promise.Promises = append(promise.Promises, DummyPromise{ "n:0", false })

	for i :=1; i < 10; i++ {
		promise.Promises = append(promise.Promises, DummyPromise{ "n:" + string(1), true })
	}

	result,_,_ := promise.Eval([]Constant{})
	equals(t, true, result)
}

func TestOrPromiseEvalAllFalse(t *testing.T) {
	promise := OrPromise{ []Promise{} }
	
	for i :=0; i < 10; i++ {
		promise.Promises = append(promise.Promises, DummyPromise{ "n:" + string(1), false })
	}

	result,_,_ := promise.Eval([]Constant{})
	equals(t, false, result)
}
