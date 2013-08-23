package llconf

import (
	"strings"
	"fmt"
)

type Promise interface {
	Desc(consts map[string][]string) string
	Eval(consts map[string][]string) bool
	SetConsts(consts map[string][]string)
}


/*
 * AndPromise
 */
type AndPromise struct {
	Promises []Promise
}

func (p AndPromise) Desc(consts map[string][]string) string {
	promises := ""
	for _,v := range(p.Promises) {
		promises += " " + v.Desc(consts)
	}
	return "(and" + promises + ")"
}

func (p AndPromise) Eval(consts map[string][]string) bool {
	for _,v := range(p.Promises) {
		if v.Eval(consts) == false {
			return false
		}
	}
	return true
}

func (p AndPromise) SetConsts(consts map[string][]string) {
	return
}

/*
 * OrPromise
 */
type OrPromise struct {
	Promises []Promise
}

func (p OrPromise) Desc(consts map[string][]string) string {
	promises := ""
	for _,v := range(p.Promises) {
		promises += " " + v.Desc(consts)
	}
	return "(or" + promises + ")"
}

func (p OrPromise) Eval(consts map[string][]string) bool {
	for _,v := range(p.Promises) {
		if v.Eval(consts) == true {
			return true
		}
	}
	return false
}

func (p OrPromise) SetConsts(consts map[string][]string) {
	return
}

/*
 * ExecPromise
 */
type ExecPromise struct {
	Constants map[string][]string
}

func (p ExecPromise) Desc(consts map[string][]string) string {


	
	command := ""
	for _,v := range(p.Constants["argument"]) {
		command += " " + v
	}

	arg := ""
	for _,v := range(p.Constants["arg"]) {
		arg += " " + v
	}

	inh_cmd := ""
	for _,v := range(consts["argument"]) {
		inh_cmd += " " + v
	}
		
	inh_arg := ""
	for _,v := range(consts["arg"]) {
		inh_arg += " " + v
	}

	return "(exec <" + command + "|" + arg +"|" + inh_cmd + "|" + inh_arg +">)"
}

func (p ExecPromise) Eval(consts map[string][]string) bool {
	return true
}

func (p ExecPromise) SetConsts(consts map[string][]string) {
	p.Constants = consts
}

/*
 * NAMED PROMISE
 */

type NamedPromise struct {
	Name string
	Promise Promise
	Constants map[string][]string
}

func (p NamedPromise) Desc(consts map[string][]string) string {
	if p.Promise != nil {
		return "(" + p.Name + " " + p.Promise.Desc(merge(p.Constants, consts)) + ")"
	} else {
		return "(" + p.Name + strings.Join(p.Constants["argument"],",")  + ">"
	}
}

func (p NamedPromise) String() string {
	return p.Desc(map[string][]string{})
}


func (p NamedPromise) Eval(map[string][]string)  bool {
	return true
}

func (p NamedPromise) SetConsts(consts map[string][]string) {
	p.Constants = consts
	fmt.Println("yyy>", strings.Join(p.Constants["argument"],","))
}

	
func merge(c1 map[string][]string, c2 map[string][]string) map[string][]string {
	result := map[string][]string{}

	for key,_ := range(c1) {
		if _,present := result[key]; present {
			result[key] = append(result[key], c1[key]...)
		} else {
			result[key] = c1[key]
		}
	}


	for key,_ := range(c2) {
		if _,present := result[key]; present {
			result[key] = append(result[key], c2[key]...)
		} else {
			result[key] = c2[key]
		}
	}

	return result
}
