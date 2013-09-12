package promise

type SetvarPromise struct {
	VarName Argument
	VarValue Argument
}

func (p SetvarPromise) Eval(arguments []Constant, logger *Logger, vars *Variables) bool {
	name  := p.VarName.GetValue(arguments, vars)
	value := p.VarValue.GetValue(arguments, vars)
	(*vars)[name] = value
	return true
}

func (p SetvarPromise) Desc(arguments []Constant) string {
	return "(setvar \"" + p.VarName.String() + "\" " + p.VarValue.String() + " )"
}

