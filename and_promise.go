package llconf

type AndPromise struct {
	Promises []Promise
}

func (p AndPromise) String() string {
	retval := ""
	for _, promise := range p.Promises {
		retval = retval + "(" + promise.String() + ")"
	}
	return retval
}

func (p AndPromise) Eval() bool {
	retval := true
	
	for _, promise := range p.Promises {
		retval = retval && promise.Eval()
	}
	return retval
}
