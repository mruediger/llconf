package parser

import (
	"testing"
	"github.com/mruediger/llconf/compiler/lexer"
)


func TestParser(t *testing.T) {
	l := lexer.Lex("test",
`(hallo welt
   (and (test "echo" "bla")
        (test "echo" "blubb")))`)
	Parse(l)
}
