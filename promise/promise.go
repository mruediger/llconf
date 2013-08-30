package promise

type Promise interface {
	Desc(arguments []Constant) string
	Eval(arguments []Constant) (bool,[]string,[]string)
}

type Argument interface {
	GetValue(arguments []Constant) string
	String() string
}
