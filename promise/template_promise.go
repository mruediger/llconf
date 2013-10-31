package promise

import (
	"os"
	"bufio"
	"strings"
	"encoding/json"
	"text/template"
)

type TemplatePromise struct {
	Arguments []Argument
}

func (t TemplatePromise) Desc(arguments []Constant) string {
	return "(template)"
}

func (t TemplatePromise) Eval(arguments []Constant, ctx *Context) bool {
	replacer := strings.NewReplacer("'", "\"")
	json_input := replacer.Replace(t.Arguments[0].GetValue(arguments, &ctx.Vars))
	template_file := t.Arguments[1].GetValue(arguments, &ctx.Vars)
	output     := t.Arguments[2].GetValue(arguments, &ctx.Vars)

	var input interface{}
	err := json.Unmarshal([]byte(json_input), &input)
	if err != nil {
		ctx.Logger.Stderr.Write([]byte(err.Error()))
		return false
	}

	tmpl, err := template.ParseFiles(template_file)
	if err != nil {
		ctx.Logger.Stderr.Write([]byte(err.Error()))
		return false
	}



	fo,err := os.Create(output)
	defer fo.Close()
	if err != nil {
		ctx.Logger.Stderr.Write([]byte(err.Error()))
		return false
	}

	bfo := bufio.NewWriter(fo)

	err = tmpl.Execute(bfo, input)

	if err != nil {
		ctx.Logger.Stderr.Write([]byte(err.Error()))
		return false
	} else {
		bfo.Flush()
		return true
	}
}
