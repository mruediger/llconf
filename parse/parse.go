package parse

import (
	"io"
	"strings"
	"strconv"
	"github.com/mruediger/llconf/promise"
)


type UnparsedPromise struct {
	Name string
	Values []UnparsedPromise
	Arguments []promise.Argument
}

func (up UnparsedPromise) String() string {
	valuesString := ""
	for _,value := range up.Values {
		valuesString += ( " " + value.String() )
	}

	constsString := ""
	for _,v := range up.Arguments {
		constsString += ( " |" + v.String() + "|" )
	}

	return "[ <" + up.Name + ">" + constsString + valuesString + " ]"
}

func (up UnparsedPromise) parse(promises map[string]promise.Promise, primary bool) promise.Promise {
	values := []promise.Promise{}

	for _,value := range up.Values {
		values = append(values, value.parse(promises, false))
	}
	
	switch up.Name {
	case "and":
		if primary {
			panic( "(and) promise not allowed as primary promise" )
		}
		return promise.AndPromise{ values }
	case "or":
		if primary {
			panic( "(or) promise not allowed as primary promise" )
		}
		return promise.OrPromise{ values }
	case "exec":
		if primary {
			panic( "(exec) promise not allowed as primary promise" )
		}
		return promise.ExecPromise{ up.Arguments }
	case "pipe":
		if primary {
			panic( "(pipe) promise not allowed as primary promise")
		}
		execs := []promise.ExecPromise{}
		for _,v := range( values ) {
			execs = append(execs, v.(promise.ExecPromise))
		}
		return promise.PipePromise{ execs }
		
	default:
		if primary {
			if len( values ) < 1 {
				panic( "need a value for NamedPromise: " + up.Name )
			}
			if len( values ) > 1 {
				panic( "to many values for NamedPromise: " + up.Name )
			}
			np := promises[up.Name].(*promise.NamedPromise)
			np.Promise = values[0]
			return np
		} else {
			if _, ok := promises[up.Name]; ok {
				np := promises[up.Name].(*promise.NamedPromise)
				return promise.NamedPromiseUsage{np, up.Arguments}
			} else {
				panic("didn't find promise (" + up.Name + ")")
			}
		}
	}
}


func readArgument( in io.RuneScanner, start rune ) promise.Argument {
	name := ""
	nameDone := false
	value := ""

	for {
		r,_,e := in.ReadRune()

		if e != nil {
			panic(e)
		}


		switch {
		case r == '"' && start == '"':
			if len(value) > 0 {
				return promise.Constant{name + ":" + value}
			} else {
				return promise.Constant{name}
			}
		case r == ']' && start == '[':
			name = strings.TrimSpace(name)
			value = strings.TrimSpace(value)

			switch name {
			case "arg":
				i,e := strconv.Atoi(strings.TrimSpace(value))
				if e != nil {
					panic(e)
				}
				return promise.ArgGetter{i}
			case "env":
				return promise.EnvGetter{value}
			default:
				panic("unknown getter type: " + name)
			}
		case r == ':':
			nameDone = true
		default:
			if !nameDone {
				name += string(r)
			} else {
				value += string(r)
			}
		}
	}
	panic("unexpected end of input")
}

func ReadPromises( in io.RuneScanner ) []UnparsedPromise {
	//skip all leading stuff till the start
	//of the first promise

	promises := []UnparsedPromise{}
	
	for {
		r,_,e := in.ReadRune()
		
		if e != nil {
			if ( e == io.EOF ) {
				return promises
			} else {
				panic(e)
			}
		}
		
		if r == '(' {
			promises = append(promises, readPromise( in ))
		}
	}
}
			
func readPromise( in io.RuneScanner ) UnparsedPromise {
	name := ""
	promises := []UnparsedPromise{}
	arguments := []promise.Argument{}

	for {
		r,_,e := in.ReadRune()
		if e != nil {
			panic(e)
		}

		switch {
		case r == '"' || r == '[':
			arguments = append(arguments, readArgument(in, r))
		case r == '(':
			promises = append(promises, readPromise(in))
		case r == ')':
			return UnparsedPromise{ strings.TrimSpace(name), promises, arguments }
		default:
			name += string(r)
		}
	}
	panic("unexpected end of input")
}





func ParsePromises( in io.RuneScanner ) map[string]promise.Promise {
	unparsed := ReadPromises( in )
	promises := map[string]promise.Promise{}

	for _,p := range( unparsed ) {
		if _,present := promises[p.Name]; present {
			panic("duplicated Promise: " + p.Name)
		} else {
			promises[p.Name] = &promise.NamedPromise { p.Name, nil }
		}
	}

	for _,p := range( unparsed ) {
		p.parse( promises, true )
	}

	return promises
}
