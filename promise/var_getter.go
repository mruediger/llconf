package promise

type VarGetter struct {
	Name string
	Vars *map[string]string
}

func (getter VarGetter) String() string {
	return "[var:" + getter.Name + "]"
}

func (getter VarGetter) GetValue(arguments []Constant) string {
	if v,present := (*getter.Vars)[getter.Name]; present {
		return v
	} else {
		return ""
	}
}
