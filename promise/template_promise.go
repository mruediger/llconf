package promise

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type TemplatePromise struct {
	JsonInput    Argument
	TemplateFile Argument
	Output       Argument
}

func (t TemplatePromise) New(children []Promise, args []Argument) (Promise, error) {
	if len(args) == 3 {
		return TemplatePromise{args[0], args[1], args[2]}, nil
	} else {
		return nil, errors.New("(template) has not enough arguments")
	}
}

func (t TemplatePromise) Desc(arguments []Constant) string {
	return fmt.Sprintf("(template in:%s temp:%s out:%s)",
		t.JsonInput,
		t.TemplateFile,
		t.Output)
}

func (t TemplatePromise) Eval(arguments []Constant, ctx *Context) bool {
	replacer := strings.NewReplacer("'", "\"")
	json_input := replacer.Replace(t.JsonInput.GetValue(arguments, &ctx.Vars))
	template_file := t.TemplateFile.GetValue(arguments, &ctx.Vars)
	output := t.Output.GetValue(arguments, &ctx.Vars)

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

	fo, err := os.Create(output)
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
