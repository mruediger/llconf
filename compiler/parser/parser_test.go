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

func TestDuplicatePromise(t *testing.T) {
	l := lexer.Lex("main.cnf", "(hallo) (hallo)")
	_,err := Parse(l)
	if err != nil {
		fmt.Println(err)
	} else {
		t.Fail()
	}
}
