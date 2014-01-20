package promise

type Variables map[string]string

type VarGetter struct {
	Name string
}

func (getter VarGetter) String() string {
	return "[var:" + getter.Name + "]"
}

func (getter VarGetter) GetValue(arguments []Constant, vars *Variables) string {
	if v, present := (*vars)[getter.Name]; present {
		return v
	} else {
		return "missing"
	}
}
