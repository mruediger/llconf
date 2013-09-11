package promise

type Variables map[string]string

type VarGetter struct {
	Name string
	Vars *Variables
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
