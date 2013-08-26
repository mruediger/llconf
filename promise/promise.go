package promise

import (
	"os"
)

type Promise interface {
	Desc(arguments []Constant) string
	Eval(arguments []Constant) bool
}


/*
 * AndPromise
 */
type AndPromise struct {
	Promises []Promise
}

func (p AndPromise) Desc(arguments []Constant) string {
	promises := ""
	for _,v := range(p.Promises) {
		promises += " " + v.Desc(arguments)
	}
	return "(and" + promises + ")"
}

func (p AndPromise) Eval(arguments []Constant) bool {
	for _,v := range(p.Promises) {
		if v.Eval(arguments) == false {
			return false
		}
	}
	return true
}


/*
 * OrPromise
 */
type OrPromise struct {
	Promises []Promise
}

func (p OrPromise) Desc(arguments []Constant) string {
	promises := ""
	for _,v := range(p.Promises) {
		promises += " " + v.Desc(arguments)
	}
	return "(or" + promises + ")"
}

func (p OrPromise) Eval(arguments []Constant) bool {
	for _,v := range(p.Promises) {
		if v.Eval(arguments) == true {
			return true
		}
	}
	return false
}

/*
 * NAMED PROMISE
 */

type NamedPromise struct {
	Name string
	Promise Promise
}

func (p NamedPromise) Desc(arguments []Constant) string {
	if p.Promise != nil {
		return "(" + p.Name + " " + p.Promise.Desc(arguments) + ")"
	} else {
		return "(" + p.Name + ">"
	}
}

func (p NamedPromise) String() string {
	return p.Desc([]Constant{})
}


func (p NamedPromise) Eval(arguments []Constant)  bool {
	return p.Promise.Eval(arguments)
}

/*
 * NAMED PROMISE USAGE
 */
type NamedPromiseUsage struct {
	Promise Promise
	Arguments []Argument
}

func (p NamedPromiseUsage) Desc(arguments []Constant) string {
	parsed_arguments := []Constant{}
	for _,argument := range(p.Arguments) {
		parsed_arguments = append(parsed_arguments, Constant{argument.GetValue(arguments)})
	}
	
	return p.Promise.Desc(append(parsed_arguments, arguments...))
}

func (p NamedPromiseUsage) Eval(arguments []Constant) bool {
	parsed_arguments := []Constant{}
	for _,argument := range(p.Arguments) {
		parsed_arguments = append(parsed_arguments, Constant{argument.GetValue(arguments)})
	}

	return p.Promise.Eval(append(parsed_arguments, arguments...))
}


type Argument interface {
	GetValue(arguments []Constant) string
	String() string
}

type Constant struct {
	Value string
}

func (constant Constant) GetValue(arguments []Constant) string {
	return constant.Value
}

func (constant Constant) String() string {
	return "constant->" + constant.Value
}

type ArgGetter struct {
	Position int
}

func (argGetter ArgGetter) GetValue(arguments []Constant) string {
	return arguments[argGetter.Position].Value
}

func (argGetter ArgGetter) String() string {
	return "arg->" + string(argGetter.Position)
}

type EnvGetter struct {
	Name string
}

func (envGetter EnvGetter) GetValue(arguments []Constant) string {
	value := os.Getenv(envGetter.Name)
	if value == "" {
		panic("didn't find environment variable " + envGetter.Name)
	}
	return value
}

func (envGetter EnvGetter) String() string {
	return "env->$" + envGetter.Name + "("+ os.Getenv(envGetter.Name) + ")"
}
