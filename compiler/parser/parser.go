package parser

import (
	"fmt"
	"errors"
	"strconv"

	"github.com/mruediger/llconf/promise"
	"github.com/mruediger/llconf/compiler/lexer"
	"github.com/mruediger/llconf/compiler/token"
)

var builtins = map[string]promise.Promise {
	"and"    : promise.AndPromise{},
	"test"   : promise.ExecPromise{Type: promise.ExecTest},
	"change" : promise.ExecPromise{Type: promise.ExecChange},
}

type UnresolvedPromise struct {
	name string
	children []UnresolvedPromise
	args []Arg
	vars []Var
	pos token.Position
}

func (p UnresolvedPromise) String() string {
	var childs string
	for _,v := range(p.children) {
		childs += v.String()
	}

	var args string
	for _,v := range(p.args) {
		args += "|"  + v.val
	}

	var vars string
	for _,v := range(p.vars) {
		vars += "[" + v.typ + ":" + v.val + "]"
	}
	return fmt.Sprintf("(%s %s %s %s)", p.name, args, p.pos ,childs);
}

func (p *UnresolvedPromise) resolvePrimary(
	unresolved map[string]UnresolvedPromise,
	builtins map[string]promise.Promise) (promise.Promise,error) {

	if len(p.children) != 1 {
		return nil, errors.New("named promise needs exactly one child, found " +
			strconv.Itoa(len(p.children)) + " " + p.pos.String())
	}

	if child,err := p.children[0].resolve(unresolved, builtins); err == nil {
		return promise.NamedPromise{Name: p.name, Promise: child}, nil
	} else {
		return nil, err
	}
}

func (p *UnresolvedPromise) resolve(
	unresolved map[string]UnresolvedPromise,
	builtins map[string]promise.Promise) (promise.Promise, error) {

	if u,present := unresolved[p.name]; present {
		return u.resolvePrimary(unresolved, builtins)
	}

	children := []promise.Promise{}
	for _,c := range p.children {
		if r,e := c.resolve(unresolved, builtins); e == nil {
			children = append(children, r)
		} else {
			return nil,e
		}
	}

	if _,present := builtins[p.name]; present {
		return builtins[p.name].New(children), nil
	}

	return nil, errors.New("couldn't find promise (" +
		p.name + ") at " + p.pos.String())
}

type Arg struct {
	val string
}

func (a *Arg) parse(l *lexer.Lexer) {
	for {
		t := l.NextToken()
		switch {
		case t.Typ == token.Argument:
			a.val = t.Val
		case t.Typ == token.RightArg:
			return
		}
	}
}

type Var struct {
	typ string
	val string
}

func (v *Var) parse(l *lexer.Lexer) {
	for {
		t := l.NextToken()
		switch {
		case t.Typ == token.GetterType:
			v.typ = t.Val
		case t.Typ == token.GetterValue:
			v.val = t.Val
		case t.Typ == token.RightGetter:
			return
		}
	}
}


func (p *UnresolvedPromise) parse(l *lexer.Lexer) error {
	for {
		t := l.NextToken()
		switch {
		case t.Typ == token.EOF:
			return nil
		case t.Typ == token.Error:
			return errors.New(t.Val)
		case t.Typ == token.LeftPromise:
			promise := UnresolvedPromise{}
			promise.pos = t.Pos
			promise.parse(l)
			p.children = append(p.children, promise)
		case t.Typ == token.PromiseName:
			p.name = t.Val
		case t.Typ == token.LeftArg:
			a := Arg{}
			a.parse(l)
			p.args = append(p.args, a)
		case t.Typ == token.LeftGetter:
			v := Var{}
			v.parse(l)
			p.vars = append(p.vars, v)
		case t.Typ == token.RightPromise:
			return nil
		}
	}
	return nil
}

func generatePromises(l *lexer.Lexer) (map[string]UnresolvedPromise, error) {

	promises := map[string]UnresolvedPromise{}

	for {
		t := l.NextToken()
		switch {
		case t.Typ == token.LeftPromise:
			p := UnresolvedPromise{}
			p.pos = t.Pos
			if err := p.parse(l); err != nil {
				return nil,err
			}
			if _,present := promises[p.name]; present {
				return nil,errors.New("found duplicate promise: " +
					p.name + " at " + p.pos.String())
			} else {
				promises[p.name] = p
			}
		case t.Typ == token.EOF:
			return promises,nil
		case t.Typ == token.Error:
			return nil,errors.New(t.Val)
		}
	}
	return nil,errors.New("unknown error")
}

func Parse(l *lexer.Lexer) (map[string]promise.Promise,error) {

	unresolved,err := generatePromises(l)
	if err != nil {
		return nil, err
	}

	resolved  := map[string]promise.Promise{}

	for k,p := range unresolved {
		if r,e := p.resolvePrimary(unresolved, builtins); e == nil {
			resolved[k] = r
		} else {
			return nil,e
		}
	}

	return resolved,nil
}
