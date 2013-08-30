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

func (p AndPromise) Eval(arguments []Constant) (bool,[]string,[]string) {
	stdout := []string{}
	stderr := []string{}
	for _,v := range(p.Promises) {
		result,o,e := v.Eval(arguments)
		stdout = append(stdout,o...)
		stderr = append(stderr,e...)
		if result == false {
			return false,stdout,stderr
		}
	}
	return true,stdout,stderr
}
