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

func (p AndPromise) Eval(arguments []Constant, logger *Logger) bool {
	for _,v := range(p.Promises) {
		result := v.Eval(arguments, logger)
		if result == false {
			return false
		}
	}
	return true
}
