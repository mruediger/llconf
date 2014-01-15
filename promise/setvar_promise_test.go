package promise

import (
	"testing"
)

func TestSetvar(t *testing.T) {
	promise := SetvarPromise{ Constant("test"), ArgGetter{0} }

	ctx := NewContext()
	promise.Eval([]Constant{Constant("foobar")}, &ctx)

	if v,ok := ctx.Vars["test"]; ok {
		equals(t, "foobar",v)
	} else {
		t.Errorf("vars[\"test\"] is undefined")
	}
}
