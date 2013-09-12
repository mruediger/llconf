package promise

type OrPromise struct {
	Promises []Promise
}

func (p OrPromise) Desc(arguments []Constant) string {
	promises := ""
	for _,v := range(p.Promises) {
		promises += " " + v.Desc(arguments)
	}
	return "(or" + promises + ")"
}

func (p OrPromise) Eval(arguments []Constant, logger *Logger, vars *Variables) bool {
	for _,v := range(p.Promises) {
		r := v.Eval(arguments, logger, vars)
		if r == true {
			return true
		}
	}
	return false
}
