package promise

import (
	"testing"
)

func TestSetvar(t *testing.T) {

	vars := Variables{}
	promise := SetvarPromise{ Constant{"test"}, ArgGetter{0} }

	promise.Eval([]Constant{Constant{"foobar"}}, &Logger{}, &vars)

	if v,ok := vars["test"]; ok {
		equals(t, "foobar",v)
	} else {
		t.Errorf("vars[\"test\"] is undefined")
	}
}
