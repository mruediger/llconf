package llconf

import (
	"io"
	"strings"
	"strconv"
)

type UnparsedPromise struct {
	Name string
	Values []UnparsedPromise
	Arguments []Argument
}

type Argument interface {
	GetValue(arguments []Constant) string
	String() string
}

type Constant struct {
	Value string
}

func (constant Constant) GetValue(arguments []Constant) string {
	return constant.Value
}

func (constant Constant) String() string {
	return "constant->" + constant.Value
}

type ArgGetter struct {
	Position int
}

func (argGetter ArgGetter) GetValue(arguments []Constant) string {
	return arguments[argGetter.Position].Value
}

func (argGetter ArgGetter) String() string {
	return string("arg->" + string(argGetter.Position))
}


// for debugging only
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

func readArgument( in io.RuneScanner ) Argument {
	name := ""
	nameDone := false
	value := ""

	for {
		r,_,e := in.ReadRune()

		if e != nil {
			panic(e)
		}


		switch r {
		case '"':
			return Constant{name}
		case ']':
			name = strings.TrimSpace(name)
			value = strings.TrimSpace(value)

			switch name {
			case "arg":
				i,e := strconv.Atoi(strings.TrimSpace(value))
				if e != nil {
					panic(e)
				}
				return ArgGetter{i}
			default:
				panic("unknown getter type: " + name)
			}
		case ':':
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
	arguments := []Argument{}

	for {
		r,_,e := in.ReadRune()
		if e != nil {
			panic(e)
		}

		switch {
		case r == '"' || r == '[':
			arguments = append(arguments, readArgument(in))
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
