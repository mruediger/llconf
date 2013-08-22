package llconf

import (
	"io"
	"strings"
)

type UnparsedPromise struct {
	Name string
	Values []UnparsedPromise
	Consts []Constant
}

type Constant struct {
	Name string
	Value string
}

// for debugging only
func (up UnparsedPromise) String() string {
	valuesString := ""
	for _,value := range up.Values {
		valuesString += ( " " + value.String() )
	}

	constsString := ""
	for _,constant := range up.Consts {
		constsString += ( " " + constant.String() )
	}

	return "[ <" + up.Name + ">" + constsString + valuesString + " ]"
}

func (c Constant) String() string {
	return "|" + c.Name + ":  " + c.Value + "|"
}

func readConstant( in io.RuneScanner ) Constant {
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
			return Constant{"arg",name}
		case ']':
			return Constant{strings.TrimSpace(name), strings.TrimSpace(value)}
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
	values := []UnparsedPromise{}
	consts := []Constant{}

	for {
		r,_,e := in.ReadRune()
		if e != nil {
			panic(e)
		}

		switch {
		case r == '"' || r == '[':
			consts = append(consts, readConstant(in))
		case r == '(':
			values = append(values, readPromise(in))
		case r == ')':
			return UnparsedPromise{ strings.TrimSpace(name), values, consts }
		default:
			name += string(r)
		}
	}
	panic("unexpected end of input")
}
