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

func (p OrPromise) Eval(arguments []Constant) (result bool, stdout []string, stderr []string) {
	for _,v := range(p.Promises) {
		r,o,e := v.Eval(arguments)
		stdout = append(stdout, o...)
		stderr = append(stderr, e...)
		if r == true {
			return true, stdout, stderr
		}
	}
	return false, stdout, stderr
}
