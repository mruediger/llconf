package promise

import "testing"

func TestNamedPromiseInvocation(t *testing.T) {
	promise := new(NamedPromiseUsage)
	promise.Promise = &NamedPromise{ "test", DummyPromise{ "content", true }}
	promise.Arguments = []Argument{ Constant{"hello"} , ArgGetter{1} }

	equals(t,
		"(test (dummy [content]constant->hello constant->foo))",
		promise.Desc([]Constant{Constant{"world"}, Constant{"foo"}}))

	equals(t, true, promise.Eval([]Constant{Constant{"world"}, Constant{"foo"}}))
}

