package llconf

import (
	"strconv"
	"strings"
)

type Promise interface {
	Desc(consts map[string][]string) string
	Eval(consts map[string][]string) bool
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


	for _,v := range(p.Constants["arg"]) {
		i,e := strconv.Atoi(strings.TrimSpace(v))
		if e != nil {
			panic(e)
		}
		command += consts["argument"][i]
	}
	
	return "(exec <" + command + ">"
}

func (p ExecPromise) Eval(consts map[string][]string) bool {
	return true
}

/*
 * NAMED PROMISE
 */

type NamedPromise struct {
	Name string
	Promise Promise
}

func (p NamedPromise) Desc(consts map[string][]string) string {
	if p.Promise != nil {
		return "(" + p.Name + " " + p.Promise.Desc(consts) + ")"
	} else {
		return "(" + p.Name + ">"
	}
}

func (p NamedPromise) String() string {
	return p.Desc(map[string][]string{})
}


func (p NamedPromise) Eval(map[string][]string)  bool {
	return true
}

/*
 * NAMED PROMISE USAGE
 */
type NamedPromiseUsage struct {
	Promise Promise
	Constants map[string][]string
}

func (p NamedPromiseUsage) Desc(consts map[string][]string) string {
	return p.Promise.Desc(merge(p.Constants, consts))
}

func (p NamedPromiseUsage) Eval(consts map[string][]string) bool {
	return p.Promise.Eval(merge(p.Constants, consts))
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
