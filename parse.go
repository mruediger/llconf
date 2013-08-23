package llconf

import (
	"io"
)


func (up UnparsedPromise) parse(promises map[string]Promise, primary bool) Promise {
	values := []Promise{}

	for _,value := range up.Values {
		values = append(values, value.parse(promises, false))
	}
	
	switch up.Name {
	case "and":
		if primary {
			panic( "(and) promise not allowed as primary promise" )
		}
		return AndPromise{ values }
	case "or":
		if primary {
			panic( "(and) promise not allowed as primary promise" )
		}
		return OrPromise{ values }
	case "exec":
		if primary {
			panic( "(and) promise not allowed as primary promise" )
		}
		return ExecPromise{ up.Consts }
	default:
		if primary {
			if len( values ) < 1 {
				panic( "need a value for NamedPromise: " + up.Name )
			}
			if len( values ) > 1 {
				panic( "to many values for NamedPromise: " + up.Name )
			}
			np := promises[up.Name].(*NamedPromise)
			np.Promise = values[0]
			return np
		} else {
			if _, ok := promises[up.Name]; ok {
				np := promises[up.Name].(*NamedPromise)
				return NamedPromiseUsage{np, up.Consts}
			} else {
				panic("didn't find promise (" + up.Name + ")")
			}
		}
	}
}



func ParsePromises( in io.RuneScanner ) map[string]Promise {
	unparsed := ReadPromises( in )
	promises := map[string]Promise{}

	for _,p := range( unparsed ) {
		if _,present := promises[p.Name]; present {
			panic("duplicated Promise: " + p.Name)
		} else {
			promises[p.Name] = &NamedPromise { p.Name, nil }
		}
	}

	for _,p := range( unparsed ) {
		p.parse( promises, true )
	}

	return promises
}
