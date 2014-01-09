package promise

type NamedPromise struct {
	Name string
	Promise Promise
}

func (p NamedPromise) New(children []Promise) Promise {
	return NamedPromise{}
}

func (p NamedPromise) Desc(arguments []Constant) string {
	if p.Promise != nil {
		return "(" + p.Name + " " + p.Promise.Desc(arguments) + ")"
	} else {
		return "(" + p.Name + " <missing promise> )"
	}
}

func (p NamedPromise) String() string {
	return p.Desc([]Constant{})
}


func (p NamedPromise) Eval(arguments []Constant, ctx *Context) bool {
	return p.Promise.Eval(arguments, ctx)
}
