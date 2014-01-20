package promise

import "fmt"

type NamedPromise struct {
	Name      string
	Promise   Promise
	Arguments []Argument
}

type NewNotSupported string

func (e NewNotSupported) Error() string {
	return "NamedPromise does not support instatiation via New()"
}

func (p NamedPromise) New(children []Promise, args []Argument) (Promise, error) {
	return nil, NewNotSupported("")
}

func (p NamedPromise) String() string {
	return p.Desc([]Constant{})
}

func (p NamedPromise) Desc(arguments []Constant) string {
	parsed_arguments := []Constant{}
	for _, argument := range p.Arguments {
		parsed_arguments = append(parsed_arguments, Constant(argument.String()))
	}

	return fmt.Sprintf("(%s %s)", p.Name, p.Promise.Desc(parsed_arguments))
}

func (p NamedPromise) Eval(arguments []Constant, ctx *Context) bool {
	parsed_arguments := []Constant{}
	for _, argument := range p.Arguments {
		parsed_arguments = append(parsed_arguments, Constant(argument.GetValue(arguments, &ctx.Vars)))
	}

	copyied_vars := Variables{}
	for k, v := range ctx.Vars {
		copyied_vars[k] = v
	}

	copyied_ctx := *ctx
	copyied_ctx.Vars = copyied_vars

	return p.Promise.Eval(parsed_arguments, &copyied_ctx)
}
