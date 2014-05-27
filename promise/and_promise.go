package promise

import "fmt"

type AndPromise struct {
	Promises []Promise
}

func (p AndPromise) New(children []Promise, args []Argument) (Promise, error) {
	if len(children) == 0 {
		return nil, fmt.Errorf("(and) needs at least 1 nested promise")
	}

	if len(args) != 0 {
		return nil, fmt.Errorf("string args are not allowed in (and) promises")
	}

	return AndPromise{children}, nil
}

func (p AndPromise) Desc(arguments []Constant) string {
	promises := ""
	for _, v := range p.Promises {
		promises += " " + v.Desc(arguments)
	}
	return "(and" + promises + ")"
}

func (p AndPromise) Eval(arguments []Constant, ctx *Context) bool {
	for _, v := range p.Promises {
		result := v.Eval(arguments, ctx)
		if result == false {
			return false
		}
	}
	return true
}
