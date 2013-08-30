package promise

// TODO rename to NamedPromiseInvocation

type NamedPromiseUsage struct {
	Promise *NamedPromise
	Arguments []Argument
}

func (p NamedPromiseUsage) Desc(arguments []Constant) string {
	parsed_arguments := []Constant{}
	for _,argument := range(p.Arguments) {
		parsed_arguments = append(parsed_arguments, Constant{argument.GetValue(arguments)})
	}
	
	return p.Promise.Desc(parsed_arguments)
}

func (p NamedPromiseUsage) Eval(arguments []Constant) (bool,[]string,[]string) {
	parsed_arguments := []Constant{}
	for _,argument := range(p.Arguments) {
		parsed_arguments = append(parsed_arguments, Constant{argument.GetValue(arguments)})
	}

	return p.Promise.Eval(parsed_arguments)
}
