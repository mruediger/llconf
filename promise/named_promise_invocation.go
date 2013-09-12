package promise

// TODO rename to NamedPromiseInvocation

type NamedPromiseUsage struct {
	Promise *NamedPromise
	Arguments []Argument
}

func (p NamedPromiseUsage) Desc(arguments []Constant) string {
	parsed_arguments := []Constant{}
	for _,argument := range(p.Arguments) {
		parsed_arguments = append(parsed_arguments, Constant{argument.String()})
	}
	
	return p.Promise.Desc(parsed_arguments)
}

func (p NamedPromiseUsage) Eval(arguments []Constant, logger *Logger, vars *Variables) bool {
	parsed_arguments := []Constant{}
	for _,argument := range(p.Arguments) {
		parsed_arguments = append(parsed_arguments, Constant{argument.GetValue(arguments, vars)})
	}

	copyied_vars := Variables{}
	for k,v := range *vars {
		copyied_vars[k] = v
	}
		
	return p.Promise.Eval(parsed_arguments, logger, &copyied_vars)
}
