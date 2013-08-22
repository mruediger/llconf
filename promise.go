package llconf

type Promise interface {
	String() string
	State() bool
}


/*
 * AndPromise
 */

type AndPromise struct {
	Promises []Promise
}

func (p AndPromise) String() string {
	promises := ""
	for _,v := range(p.Promises) {
		promises += " " + v.String()
	}
	return "(and" + promises + ")"
}

func (p AndPromise) State() bool {
	for _,v := range(p.Promises) {
		if v.State() == false {
			return false
		}
	}
	return true
}

/*
 * OrPromise
 */
type OrPromise struct {
	Promises []Promise
}

func (p OrPromise) String() string {
	promises := ""
	for _,v := range(p.Promises) {
		promises += " " + v.String()
	}
	return "(or" + promises + ")"
}

func (p OrPromise) State() bool {
	for _,v := range(p.Promises) {
		if v.State() == true {
			return true
		}
	}
	return false
}

/*
 * ExecPromise
 */
type ExecPromise struct {
	Command string
}

func (p ExecPromise) String() string {
	return "(exec " + p.Command + ")"
}

func (p ExecPromise) State() bool {
	return true
}

/*
 * NAMED PROMISE
 */

type NamedPromise struct {
	Name string
	Promise *Promise 
}

func (p NamedPromise) String() string {
	if p.Promise != nil {
		return "(" + p.Name + " " + (*p.Promise).String() + ")"
	} else {
		return "(" + p.Name + ")"
	}
}

func (p NamedPromise) State() bool {
	return true
}


	
