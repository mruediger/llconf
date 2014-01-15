package promise

type Constant string

func (constant Constant) GetValue(arguments []Constant, vars *Variables) string {
	return string(constant)
}

func (constant Constant) String() string {
	return "constant->" + string(constant)
}
