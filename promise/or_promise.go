package promise

import "fmt"

type OrPromise struct {
	Promises []Promise
}

func (p OrPromise) New(children []Promise, args []Argument) (Promise, error) {
	if len(children) == 0 {
		return nil, fmt.Errorf("(and) needs at least 1 nested promise")
	}

	if len(args) != 0 {
		return nil, fmt.Errorf("string args are not allowed in (and) promises")
	}

	return OrPromise{children}, nil
}

func (p OrPromise) Desc(arguments []Constant) string {
	promises := ""
	for _, v := range p.Promises {
		promises += " " + v.Desc(arguments)
	}
	return "(or" + promises + ")"
}

func (p OrPromise) Eval(arguments []Constant, ctx *Context) bool {
	for _, v := range p.Promises {
		r := v.Eval(arguments, ctx)
		if r == true {
			return true
		}
	}
	return false
}
