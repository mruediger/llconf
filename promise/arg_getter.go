package promise

type ArgGetter struct {
	Position int
}

func (argGetter ArgGetter) GetValue(arguments []Constant) string {
	if len(arguments) <= argGetter.Position {
		return ""
	}
	return arguments[argGetter.Position].Value
}

func (argGetter ArgGetter) String() string {
	return "arg->" + string(argGetter.Position)
}
