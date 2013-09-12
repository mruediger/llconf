package promise

type Constant struct {
	Value string
}

func (constant Constant) GetValue(arguments []Constant, vars *Variables) string {
	return constant.Value
}

func (constant Constant) String() string {
	return "constant->" + constant.Value
}
