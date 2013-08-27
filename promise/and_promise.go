package promise

type AndPromise struct {
	Promises []Promise
}

func (p AndPromise) Desc(arguments []Constant) string {
	promises := ""
	for _,v := range(p.Promises) {
		promises += " " + v.Desc(arguments)
	}
	return "(and" + promises + ")"
}

func (p AndPromise) Eval(arguments []Constant) bool {
	for _,v := range(p.Promises) {
		if v.Eval(arguments) == false {
			return false
		}
	}
	return true
}
