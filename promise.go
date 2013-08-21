package llconf

type Promise interface {
	Eval() bool
	String() string
	
}

type OrPromise struct {
	Promises []Promise
}

type ExecPromise struct {
	Command string
}

type NamedPromise struct {
	Name string
	Promise
}

