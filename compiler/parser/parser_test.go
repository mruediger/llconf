package parser

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/mruediger/llconf/promise"
)

func TestParser(t *testing.T) {
	_, err := Parse([]Input{{"main.cnf",
		`(hallo welt
 (and (test "echo" "foo bar")
        (test "echo" "blubb")
        (change "bla" [var:blubb])))`}})
	if err != nil {
		t.Errorf("TestParse: " + err.Error())
	}
}

func TestMultiplePromises(t *testing.T) {
	_, err := Parse([]Input{{"main.cnf", "(hallo (test)) (welt (test))"}})
	if err != nil {
		t.Errorf("MultiplePromises: " + err.Error())
	}
}

func TestUsePromise(t *testing.T) {
	p, err := Parse([]Input{{"main.cnf",
		`(hallo (welt))
 (welt (and (test "echo" "foo") (test "echo" "bar")))`}})

	if err == nil {

		w1 := p["hallo"].(promise.NamedPromise).Promise.(promise.NamedPromise)
		w2 := p["welt"].(promise.NamedPromise)

		w1.Name = "w1"
		w2.Name = "w2"

		if reflect.DeepEqual(w1, w2) {
			t.Errorf("multiple promises are the same object")
		}
	} else {
		fmt.Println(err)
	}
}

func TestUseVars(t *testing.T) {
	p, err := Parse([]Input{{"main.cnf",
		`(t1 (hallo "a"))
 (t2 (hallo "b"))
 (hallo (test "echo" [arg:0]))`}})

	if err == nil {
		if !strings.Contains(p["t1"].(promise.NamedPromise).String(), "constant->a") {
			t.Errorf("TestUseVars: couldn't find value")
		}
		if !strings.Contains(p["t2"].(promise.NamedPromise).String(), "constant->b") {
			t.Errorf("TestUseVars: couldn't find value")
		}
	} else {
		t.Errorf("TestUseVars: " + err.Error())
	}
}

func TestGetter(t *testing.T) {
	p, err := Parse([]Input{{"main.cnf", "(hallo (test \"echo\" [env:home] [var:test]))"}})
	if err != nil {
		t.Errorf("TestGetter: %s", err)
	} else {
		exec := p["hallo"].(promise.NamedPromise).Promise.(promise.ExecPromise)

		e_name := exec.Arguments[1].(promise.EnvGetter).Name
		v_name := exec.Arguments[2].(promise.VarGetter).Name

		if e_name != "home" {
			t.Errorf("TestGetter: env name not found")
		}
		if v_name != "test" {
			t.Errorf("TestGetter: var name not found")
		}
	}
}

func TestUnknownPromise(t *testing.T) {
	_, err := Parse([]Input{{"main.cnf", "(hallo (welt))"}})
	if err == nil {
		t.Errorf("TestUnknownPromise: expected exception")
	}
}

func TestDuplicatePromise(t *testing.T) {
	_, err := Parse([]Input{{"main.cnf", "(hallo) (hallo)"}})
	if err == nil {
		t.Errorf("TestDuplicatePromise: expected exception")
	}
}

func TestMultipleInputs(t *testing.T) {
	_, err := Parse([]Input{
		{"main.cnf", "(hallo (welt))"},
		{"welt.cnf", "(welt (test \"echo\" \"foo\"))"},
	})

	if err != nil {
		t.Errorf("TestMultipleInputs: " + err.Error())
	}
}
