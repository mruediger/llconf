package lexer

import (
	"fmt"
	"testing"

	"github.com/mruediger/llconf/compiler/token"
)

type test struct {
	name string
	input string
	output []testToken
}

type testToken struct {
	typ token.Type
	pos int
	val string
}

func (t testToken) String() string {
	return fmt.Sprintf("{%s, %d, %q}", t.typ, t.pos, t.val)
}

func equals(got, expected []testToken) bool {
	if len(got) != len(expected) {
		return false
	}

	for k := range got {
		if got[k].typ != expected[k].typ {
			return false
		}
		if got[k].pos != expected[k].pos {
			return false
		}
		if got[k].val != expected[k].val {
			return false
		}
	}


	return true;
}

var	tests = []test{
	{"basic", "(hello world  )", []testToken{
		{token.LeftPromise, 0, "("},
		{token.PromiseName, 1, "hello world"},
		{token.RightPromise, 14, ")"},
		{token.EOF, 15, ""}}},
	{"unclosed promise", "(hello world", []testToken{
		{token.LeftPromise, 0, "("},
		{token.Error, 1, "unexpected eof in promise"}}},
	{"nested promise", "(hello ( world ))", []testToken{
		{token.LeftPromise, 0, "("},
		{token.PromiseName, 1, "hello"},
		{token.LeftPromise, 7, "("},
		{token.PromiseName, 8, "world"},
		{token.RightPromise, 15, ")"},
		{token.RightPromise, 16, ")"},
		{token.EOF, 17, ""}}},
	{"mixed promise", "(hello ( world ) \"bla\")", []testToken{
		{token.LeftPromise, 0, "("},
		{token.PromiseName, 1, "hello"},
		{token.LeftPromise, 7, "("},
		{token.PromiseName, 8, "world"},
		{token.RightPromise, 15, ")"},
		{token.LeftArg, 16, "\""},
		{token.Argument, 18, "bla"},
		{token.RightArg, 21, "\""},
		{token.RightPromise, 22, ")"},
		{token.EOF, 23, ""}}},

	{"arg", "(test \"bla\"  )", []testToken{
		{token.LeftPromise, 0, "("},
		{token.PromiseName, 1, "test"},
		{token.LeftArg, 6, "\""},
		{token.Argument, 7, "bla"},
		{token.RightArg, 10, "\""},
		{token.RightPromise, 11, ")"},
		{token.EOF, 14, ""}}},
	{"two args", "(test \"bla\" \"blubb\" )", []testToken{
		{token.LeftPromise, 0, "("},
		{token.PromiseName, 1, "test"},
		{token.LeftArg, 6, "\""},
		{token.Argument, 7, "bla"},
		{token.RightArg, 10, "\""},
		{token.LeftArg, 11, "\""},
		{token.Argument, 13, "blubb"},
		{token.RightArg, 18, "\""},
		{token.RightPromise, 19, ")"},
		{token.EOF, 21, ""}}},

	{"unclosed arg", "(test \"bla)", []testToken{
		{token.LeftPromise, 0, "("},
		{token.PromiseName, 1, "test"},
		{token.LeftArg, 6, "\""},
		{token.Error, 7, "unexpected eof in argument"}}},

	{"getter", "(test [var:test])", []testToken{
		{token.LeftPromise, 0, "("},
		{token.PromiseName, 1, "test"},
		{token.LeftGetter, 6, "["},
		{token.GetterType, 7, "var"},
		{token.GetterSeparator, 10, ":"},
		{token.GetterValue, 11, "test"},
		{token.RightGetter, 15, "]"},
		{token.RightPromise, 16, ")"},
		{token.EOF, 17, ""},
	}},
}



func TestLexer(t *testing.T) {

	for _, test := range tests {
		output := runTest(test)

		if !equals(output, test.output) {
			t.Errorf("%s: got\n%v\nexpected\n%v\n", test.name, output, test.output)
		}

		fmt.Println(output)
	}
}

func runTest(test test) []testToken {

	l := Lex(test.name, test.input)

	output := []testToken{}

	for {
		t:= <- l.tokens
		output = append(output, testToken{t.Typ, t.Pos.Start, t.Val})
		if (t.Typ == token.Error || t.Typ == token.EOF) {
			break
		}
	}
	return output
}
