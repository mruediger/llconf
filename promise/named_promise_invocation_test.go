package promise

import "testing"

func TestNamedPromiseInvocation(t *testing.T) {
	promise := new(NamedPromiseUsage)
	promise.Promise = &NamedPromise{ "test", DummyPromise{ "content", true }}
	promise.Arguments = []Argument{ Constant{"hello"} , ArgGetter{1} }

	result := promise.Eval([]Constant{Constant{"world"}, Constant{"foo"}}, &Context{})
	equals(t, true, result)
}

func TestVarCopying(t *testing.T) {

	promise := new(NamedPromiseUsage)
	promise.Promise = &NamedPromise{"test", SetvarPromise{Constant{"test"}, Constant{"blafasel"}}}
	promise.Arguments = []Argument{}

	ctx := NewContext()
	ctx.Vars["test"] = "hello world"

	promise.Eval([]Constant{}, &ctx)

	equals(t, "hello world", ctx.Vars["test"])
}
