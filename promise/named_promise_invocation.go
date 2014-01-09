package promise

// TODO rename to NamedPromiseInvocation

type NamedPromiseUsage struct {
	Promise *NamedPromise
	Arguments []Argument
}

func (p NamedPromiseUsage) New(children []Promise) Promise {
	return NamedPromise{}
}

func (p NamedPromiseUsage) Desc(arguments []Constant) string {
	parsed_arguments := []Constant{}
	for _,argument := range(p.Arguments) {
		parsed_arguments = append(parsed_arguments, Constant{argument.String()})
	}

	return p.Promise.Desc(parsed_arguments)
}

func (p NamedPromiseUsage) Eval(arguments []Constant, ctx *Context) bool {
	parsed_arguments := []Constant{}
	for _,argument := range(p.Arguments) {
		parsed_arguments = append(parsed_arguments, Constant{argument.GetValue(arguments, &ctx.Vars)})
	}

	copyied_vars := Variables{}
	for k,v := range ctx.Vars {
		copyied_vars[k] = v
	}

	copyied_ctx := *ctx
	copyied_ctx.Vars = copyied_vars



	return p.Promise.Eval(parsed_arguments, &copyied_ctx)
}
