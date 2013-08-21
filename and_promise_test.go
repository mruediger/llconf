package llconf


import (
	"testing"
)

func TestAndPromise(t *testing.T) {

	promises := []Promise {
		DummyPromise{ "test_a", true },
		DummyPromise{ "test_b", false },
		DummyPromise{ "test_c", true }}

	p := AndPromise{ promises }

	var and_string = p.String()
	var and_string_want = "(test_a)(test_b)(test_c)"

	if and_string != and_string_want {
		t.Errorf("p.String() == %q, want %q", and_string, and_string_want)
	}

	var and_value = p.Eval()
	var and_value_want = false

	if and_value != and_value_want {
		t.Errorf("p.Eval == %q, want %q", and_value, and_value_want)
	}
}
