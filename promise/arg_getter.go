package promise

type ArgGetter struct {
	Position int
}

func (argGetter ArgGetter) GetValue(arguments []Constant) string {
	return arguments[argGetter.Position].Value
}

func (argGetter ArgGetter) String() string {
	return "arg->" + string(argGetter.Position)
}
