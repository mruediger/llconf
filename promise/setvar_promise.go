package promise

type SetvarPromise struct {
	VarName Argument
	VarValue Argument
}

func (p SetvarPromise) New(children []Promise) Promise {
	return SetvarPromise{}
}

func (p SetvarPromise) Eval(arguments []Constant, ctx *Context) bool {
	name  := p.VarName.GetValue(arguments, &ctx.Vars)
	value := p.VarValue.GetValue(arguments, &ctx.Vars)
	ctx.Vars[name] = value
	return true
}

func (p SetvarPromise) Desc(arguments []Constant) string {
	return "(setvar \"" + p.VarName.String() + "\" " + p.VarValue.String() + " )"
}
