package promise

type SetvarPromise struct {
	VarName  Argument
	VarValue Argument
}

type SetvarError string

func (e SetvarError) Error() string {
	return string(e)
}

const (
	ArgumentCount = SetvarError("(setvar) needs 2 arguments")
)

func (p SetvarPromise) New(children []Promise, args []Argument) (Promise, error) {

	if len(args) != 2 {
		return nil, ArgumentCount
	}

	setvar := SetvarPromise{}
	setvar.VarName = args[0]
	setvar.VarValue = args[1]

	return setvar, nil
}

func (p SetvarPromise) Eval(arguments []Constant, ctx *Context, stack string) bool {
	name := p.VarName.GetValue(arguments, &ctx.Vars)
	value := p.VarValue.GetValue(arguments, &ctx.Vars)
	ctx.Vars[name] = value
	return true
}

func (p SetvarPromise) Desc(arguments []Constant) string {
	return "(setvar \"" + p.VarName.String() + "\" " + p.VarValue.String() + " )"
}
