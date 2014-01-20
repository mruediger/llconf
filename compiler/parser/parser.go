package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/mruediger/llconf/compiler/lexer"
	"github.com/mruediger/llconf/compiler/token"
	"github.com/mruediger/llconf/promise"
)

var builtins = map[string]promise.Promise{
	"or":   promise.OrPromise{},
	"and":  promise.AndPromise{},
	"test": promise.ExecPromise{Type: promise.ExecTest},

	"indir":    promise.InDir{},
	"setenv":   promise.SetEnv{},
	"change":   promise.ExecPromise{Type: promise.ExecChange},
	"pipe":     promise.PipePromise{},
	"setvar":   promise.SetvarPromise{},
	"readvar":  promise.ReadvarPromise{},
	"template": promise.TemplatePromise{},
	"restart":  promise.RestartPromise{},
}

type UnresolvedPromise struct {
	name     string
	children []UnresolvedPromise
	args     []promise.Argument
	pos      token.Position
}

type Input struct {
	File   string
	String string
}

func (p UnresolvedPromise) String() string {
	var childs string
	for _, v := range p.children {
		childs += v.String()
	}

	var args string
	for _, v := range p.args {
		args += "|" + v.String()
	}
	return fmt.Sprintf("(%s %s %s %s)", p.name, args, p.pos, childs)
}

func (p *UnresolvedPromise) resolvePrimary(
	unresolved Tree,
	builtins map[string]promise.Promise) (promise.NamedPromise, error) {

	if len(p.children) != 1 {
		return promise.NamedPromise{}, errors.New("named promise needs exactly one child, found " +
			strconv.Itoa(len(p.children)) + " " + p.pos.String())
	}

	if child, err := p.children[0].resolve(unresolved, builtins); err == nil {
		return promise.NamedPromise{Name: p.name, Promise: child}, nil
	} else {
		return promise.NamedPromise{}, err
	}
}

func (p *UnresolvedPromise) resolve(
	unresolved Tree,
	builtins map[string]promise.Promise) (promise.Promise, error) {

	if u, present := unresolved[p.name]; present {
		t, e := u.resolvePrimary(unresolved, builtins)
		t.Arguments = p.args
		return t, e
	}

	children := []promise.Promise{}
	for _, c := range p.children {
		if r, e := c.resolve(unresolved, builtins); e == nil {
			children = append(children, r)
		} else {
			return nil, e
		}
	}

	if _, present := builtins[p.name]; present {
		if promise, err := builtins[p.name].New(children, p.args); err == nil {
			return promise, nil
		} else {
			return nil, errors.New(err.Error() + " at " + p.pos.String())
		}
	}

	return nil, errors.New("couldn't find promise (" +
		p.name + ") at " + p.pos.String())
}

func parseGetter(l *lexer.Lexer) (promise.Argument, error) {
	var typ string
	var getter promise.Argument
	for {
		t := l.NextToken()
		switch {
		case t.Typ == token.Error:
			return nil, errors.New(t.Val + " " + t.Pos.String())
		case t.Typ == token.GetterType:
			if t.Val == "join" {
				return parseJoiner(l)
			} else {
				typ = t.Val
			}
		case t.Typ == token.GetterValue:
			switch typ {
			case "arg":
				if i, e := strconv.Atoi(t.Val); e == nil {
					getter = promise.ArgGetter{i}
				} else {
					return nil, e
				}
			case "env":
				getter = promise.EnvGetter{t.Val}
			case "var":
				getter = promise.VarGetter{t.Val}
			default:
				return nil, fmt.Errorf("unknown getter type: %q", t.Val)
			}
		case t.Typ == token.RightGetter:
			return getter, nil
		}
	}
}

func parseJoiner(l *lexer.Lexer) (promise.Argument, error) {

	joiner := promise.JoinArgument{}

	for {
		t := l.NextToken()
		switch t.Typ {
		case token.LeftArg:
			if arg, err := parseArg(l); err == nil {
				joiner.Args = append(joiner.Args, arg)
			} else {
				return nil, err
			}
		case token.LeftGetter:
			if get, err := parseGetter(l); err == nil {
				joiner.Args = append(joiner.Args, get)
			} else {
				return nil, err
			}
		case token.RightGetter:
			return joiner, nil
		default:
			return nil, fmt.Errorf("unexpected token in joiner: %q in %s", t.Val, t.Pos.String())
		}
	}
}

func parseArg(l *lexer.Lexer) (promise.Argument, error) {
	var arg promise.Constant
	for {
		t := l.NextToken()
		switch {
		case t.Typ == token.Error:
			return nil, errors.New(t.Val + " " + t.Pos.String())
		case t.Typ == token.Argument:
			arg = promise.Constant(t.Val)
		case t.Typ == token.RightArg:
			return arg, nil
		default:
			return nil, errors.New("unknown getter type: " + t.Val)
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
			return errors.New(t.Val + " " + t.Pos.String())
		case t.Typ == token.LeftPromise:
			promise := UnresolvedPromise{}
			promise.pos = t.Pos
			if err := promise.parse(l); err == nil {
				p.children = append(p.children, promise)
			} else {
				return err
			}
		case t.Typ == token.PromiseName:
			p.name = t.Val
		case t.Typ == token.LeftArg:
			if a, e := parseArg(l); e == nil {
				p.args = append(p.args, a)
			} else {
				return e
			}
		case t.Typ == token.LeftGetter:
			if g, e := parseGetter(l); e == nil {
				p.args = append(p.args, g)
			} else {
				return e
			}
		case t.Typ == token.RightPromise:
			return nil
		}
	}
	return nil
}

type Tree map[string]UnresolvedPromise

func (tree Tree) generatePromises(l *lexer.Lexer) error {
	for {
		t := l.NextToken()
		switch {
		case t.Typ == token.LeftPromise:
			p := UnresolvedPromise{}
			p.pos = t.Pos
			if err := p.parse(l); err != nil {
				return err
			}
			if _, present := tree[p.name]; present {
				return errors.New("found duplicate promise: " +
					p.name + " at " + p.pos.String())
			} else {
				tree[p.name] = p
			}
		case t.Typ == token.EOF:
			return nil
		case t.Typ == token.Error:
			return errors.New(t.Val + " " + t.Pos.String())
		}
	}
	return errors.New("unknown error")
}

func Parse(inputs []Input) (map[string]promise.Promise, error) {
	unresolved := Tree{}

	for _, input := range inputs {
		l := lexer.Lex(input.File, input.String)
		err := unresolved.generatePromises(l)
		if err != nil {
			return nil, err
		}
	}

	resolved := map[string]promise.Promise{}

	for k, p := range unresolved {
		if r, e := p.resolvePrimary(unresolved, builtins); e == nil {
			resolved[k] = r
		} else {
			return nil, e
		}
	}

	return resolved, nil
}
