package token

import "fmt"

type Position struct {
	File  string
	Line  int
	Start int
	End   int
}

func (p Position) String() string {
	return fmt.Sprintf("%q line %d",
		p.File,
		p.Line)
}

type Type int

const (
	Error Type = iota
	Comment
	LeftPromise
	RightPromise
	PromiseName
	LeftArg
	Argument
	RightArg
	LeftGetter
	GetterType
	GetterSeparator
	GetterValue
	RightGetter
	EOF
)

var tokenNames = [...]string{
	Error:           "Error",
	LeftPromise:     "LeftPromise",
	RightPromise:    "RightPromise",
	PromiseName:     "PromiseName",
	LeftArg:         "LeftArg",
	Argument:        "Argument",
	RightArg:        "RightArg",
	LeftGetter:      "LeftGetter",
	GetterType:      "GetterType",
	GetterSeparator: "GetterSeparator",
	GetterValue:     "GetterValue",
	RightGetter:     "RightGetter",
	EOF:             "EOF",
}

func (t Type) String() string {
	return tokenNames[t]
}

type Token struct {
	Typ Type
	Pos Position
	Val string
}

func (t Token) String() string {
	return fmt.Sprintf("{%s, %s, %q}",
		t.Typ,
		t.Pos,
		t.Val)

}
