package parser

import (
	"fmt"
	"testing"
	"github.com/mruediger/llconf/compiler/lexer"
)


func TestParser(t *testing.T) {
	l := lexer.Lex("main.cnf",
`(hallo welt
 (and (test "echo" "foo bar")
        (test "echo" "blubb")
        (change "bla" [var:blubb])))`)
	if p,err := Parse(l); err == nil {
		fmt.Println(p["hallo welt"])
	} else {
		t.Fail()
	}
}

func TestMultiplePromises(t *testing.T) {
	l := lexer.Lex("main.cnf", "(hallo (test)) (welt (test))")
	if p,err := Parse(l); err == nil {
		fmt.Println(p)
	} else {
		t.Errorf("MultiplePromises: " + err.Error())
	}
}

func TestUsePromise(t *testing.T) {
	l := lexer.Lex("main.cnf",
`(hallo (welt))
 (welt (and (test "echo" "foo") (test "echo" "bar")))`)

	if p,err := Parse(l); err == nil {
		fmt.Println(p)
	} else {
		fmt.Println(err)
	}
}

func TestUnknownPromise(t *testing.T) {
	l := lexer.Lex("main.cnf", "(hallo (welt))")
	_,err := Parse(l)
	if err == nil {
		t.Errorf("TestDuplicatePromise: expected exception")
	}
}


func TestDuplicatePromise(t *testing.T) {
	l := lexer.Lex("main.cnf", "(hallo) (hallo)")
	_,err := Parse(l)
	if err == nil {
		t.Errorf("TestDuplicatePromise: expected exception")
	}
}
