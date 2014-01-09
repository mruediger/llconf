package promise

type AndPromise struct {
	Promises []Promise
}

func (p AndPromise) New(children []Promise) Promise {
	return AndPromise{children}
}

func (p AndPromise) Desc(arguments []Constant) string {
	promises := ""
	for _,v := range(p.Promises) {
		promises += " " + v.Desc(arguments)
	}
	return "(and" + promises + ")"
}

func (p AndPromise) Eval(arguments []Constant, ctx *Context) bool {
	for _,v := range(p.Promises) {
		result := v.Eval(arguments, ctx)
		if result == false {
			return false
		}
	}
	return true
}
