package main

import (
	"io"
	"os"
	"bufio"
	"fmt"
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

func (up UnparsedPromise) String() string {
	valuesString := ""
	for _,value := range up.Values {
		valuesString += (" " + value.String() )
	}

	constsString := ""
	for _,constant := range up.Consts {
		constsString += ( " " + constant.String() )
	}

	return "[ <" + up.Name + ">" + constsString + valuesString + " ]"
}

func (c Constant) String() string {
	return "|" + c.Name + ": " + c.Value + "|"
}

func readConstant( in io.RuneScanner ) Constant {

	name := ""
	nameDone := false
	value := ""

	for {
		r,_,e := in.ReadRune()

		if e != nil {
			if e == io.EOF {
				fmt.Println("unexpected EOF in constant")
			} else {
				fmt.Println("unexpected error: ",e)
			}
			break
		}
		switch r {
		case '"':
			return Constant{"arg",name}
		case ']':
			return Constant{name, value}
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
	return Constant{"error","error"}

}

func readPromise( in io.RuneScanner ) UnparsedPromise {

	name := ""
	values := make([]UnparsedPromise,0)
	consts := make([]Constant, 0)

	for {
		r,_,e := in.ReadRune()

		if e != nil {
			if e == io.EOF {
				fmt.Println("unexpected EOF in promise " + name)
			} else {
				fmt.Println("unexpected error: ",e)
			}
			break
		}
		
		switch {
		case r == '"' || r == '[':
			consts = append(consts, readConstant(in))
		case r == '(':
			values = append(values, readPromise(in))
		case r == ')':
			return UnparsedPromise{strings.TrimSpace(name), values, consts}
		default:
			name += string(r)
		}
	}
	return UnparsedPromise{"undef", values, consts}
}

func ReadLLC( in io.RuneScanner ) {
	for {
		r, _, e := in.ReadRune()
		if e != nil {
			break
		}

		if r == '(' {
			prepromise := readPromise(in)
			fmt.Println(prepromise)
		}
	}
}

func main() {
	in := os.Stdin
	bufin := bufio.NewReader( in )
	ReadLLC( bufin )
}
