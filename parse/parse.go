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

func (up UnparsedPromise) parse(promises map[string]promise.Promise, primary bool) (promise.Promise, error) {
	values := []promise.Promise{}

	for _,value := range up.Values {
		value, err := value.parse(promises, false)
		if err == nil {
			values = append(values, value)
		} else {
			return nil, err
		}
	}
	
	switch up.Name {
	case "and":
		if primary {
			return nil, IllegalPromisePosition{"and"}
		}
		return promise.AndPromise{ values },nil
	case "or":
		if primary {
			return nil, IllegalPromisePosition{"or"}
		}
		return promise.OrPromise{ values },nil
	case "exec":
		if primary {
			return nil, IllegalPromisePosition{"exec"}
		}
		return promise.ExecPromise{ up.Arguments },nil
	case "pipe":
		if primary {
			return nil, IllegalPromisePosition{"pipe"}
		}
		execs := []promise.ExecPromise{}
		for _,v := range( values ) {
			execs = append(execs, v.(promise.ExecPromise))
		}
		return promise.PipePromise{ execs }, nil
		
	default:
		if primary {
			if len( values ) != 1 {
				return nil,NamedPromiseArgc{ len(values), up.Name }
			}
			np := promises[up.Name].(*promise.NamedPromise)
			np.Promise = values[0]
			return np,nil
		} else {
			if _, ok := promises[up.Name]; ok {
				np := promises[up.Name].(*promise.NamedPromise)
				return promise.NamedPromiseUsage{np, up.Arguments}, nil
			} else {
				return nil,MissingPromise{ up.Name }
			}
		}
	}
}


func readArgument( in io.RuneReader, start rune ) (promise.Argument, error) {
	name := ""
	nameDone := false
	value := ""

	for {
		r,_,e := in.ReadRune()

		if e != nil {
			return nil,e
		}


		switch {
		case r == '"' && start == '"':
			if len(value) > 0 {
				return promise.Constant{name + ":" + value},nil
			} else {
				return promise.Constant{name},nil
			}
		case r == ']' && start == '[':
			name = strings.TrimSpace(name)
			value = strings.TrimSpace(value)

			switch name {
			case "arg":
				i,e := strconv.Atoi(strings.TrimSpace(value))
				if e != nil {
					return nil,e
				}
				return promise.ArgGetter{i},nil
			case "env":
				return promise.EnvGetter{value},nil
			default:
				return nil, UnknownGetterType{name}
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
	return nil, UnexpectedEOF{}
}

func ReadPromises( in io.RuneReader ) ([]UnparsedPromise,error) {
	//skip all leading stuff till the start
	//of the first promise

	promises := []UnparsedPromise{}
	
	for {
		r,_,e := in.ReadRune()
		
		if e != nil {
			if ( e == io.EOF ) {
				return promises,nil
			} else {
				return []UnparsedPromise{},e
			}
		}
		
		if r == '(' {
			promise,err := readPromise( in )
			if err == nil{
				promises = append(promises, promise)
			} else {
				return nil,err
			}
		}
	}
}
			
func readPromise( in io.RuneReader ) (UnparsedPromise,error) {
	name := ""
	promises := []UnparsedPromise{}
	arguments := []promise.Argument{}

	for {
		r,_,e := in.ReadRune()
		if e != nil {
			return UnparsedPromise{},e
		}

		switch {
		case r == '"' || r == '[':
			argument, err := readArgument(in,r)
			if err == nil {
				arguments = append(arguments, argument)
			} else {
				return UnparsedPromise{},e
			}
		case r == '(':
			promise,err := readPromise(in)
			if err == nil {
				promises = append(promises, promise)
			} else {
				return UnparsedPromise{},e
			}
		case r == ')':
			return UnparsedPromise{ strings.TrimSpace(name), promises, arguments }, nil
		default:
			name += string(r)
		}
	}
	return UnparsedPromise{}, UnexpectedEOF{}
}





func ParsePromises( in io.RuneReader ) (map[string]promise.Promise,error) {
	promises := map[string]promise.Promise{}
	
	unparsed,err := ReadPromises( in )
	if err != nil {
		return promises,err
	}
	

	for _,p := range( unparsed ) {
		if _,present := promises[p.Name]; present {
			return map[string]promise.Promise{}, DuplicatePromise{ p.Name }
		} else {
			promises[p.Name] = &promise.NamedPromise { p.Name, nil }
		}
	}

	for _,p := range( unparsed ) {
		_,err := p.parse( promises, true )
		if err != nil {
			return promises,err
		}
	}

	return promises,nil
}
