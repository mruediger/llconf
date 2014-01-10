package promise

import (
	"testing"
	"strings"
)

type DummyPromise struct {
	StringValue string
	EvalValue bool
}

func (p DummyPromise) New(children []Promise, args []Argument) (Promise,error) {
	return DummyPromise{},nil
}

func (p DummyPromise) Desc(arguments []Constant) string {
	var args []string

	for _,argument := range( arguments ) {
		args = append(args, argument.String())
	}
	return "(dummy [" + p.StringValue + "]" + strings.Join(args," ") + ")"
}

func (p DummyPromise) Eval(arguments []Constant, ctx *Context) bool {
	return p.EvalValue
}

func equals(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("error: wantet %q, got %q", a, b)
	}
}
