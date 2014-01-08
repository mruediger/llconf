package parser

import (
	"fmt"

	"github.com/mruediger/llconf/promise"
	"github.com/mruediger/llconf/compiler/lexer"
	"github.com/mruediger/llconf/compiler/token"
)

type UnresolvedPromise struct {
	name string
	children []UnresolvedPromise
	args string
}

func (p *UnresolvedPromise) resolve() promise.Promise {

	switch p.name {
	default:
		return promise.NamedPromise{Name: p.name}
	}

}

func (p *UnresolvedPromise) parse(l *lexer.Lexer) {
	for {
		switch t := l.NextToken(); {
		case t.Typ == token.EOF:
			return
		case t.Typ == token.LeftPromise:
			promise := UnresolvedPromise{}
			promise.parse(l)
			p.children = append(p.children, promise)
		case t.Typ == token.PromiseName:
			p.name = t.Val
		case t.Typ == token.RightPromise:
			return
		}
	}
}

func Parse(l *lexer.Lexer) {

	for l.NextToken().Typ != token.LeftPromise {}

	root := UnresolvedPromise{}
	root.parse(l)

	fmt.Println(root.resolve())

	fmt.Println(root)
}
