package promise

import (
	"os"
	"bufio"
	"encoding/json"
	"text/template"
)

type TemplatePromise struct {
	Arguments []Argument
}

func (t TemplatePromise) Desc(arguments []Constant) string {
	
	return "hello"
}

func (t TemplatePromise) Eval(arguments []Constant, logger *Logger, vars *Variables) bool {
	json_input := t.Arguments[0].GetValue(arguments,vars)
	template_file := t.Arguments[1].GetValue(arguments,vars)
	output     := t.Arguments[2].GetValue(arguments,vars)

	var input interface{}
	err := json.Unmarshal([]byte(json_input), &input)
	if err != nil {
		logger.Stderr.Write([]byte(err.Error()))
		return false
	}

	tmpl, err := template.ParseFiles(template_file)
	if err != nil {
		logger.Stderr.Write([]byte(err.Error()))
		return false
	}



	fo,err := os.Create(output)
	defer fo.Close()
	if err != nil {
		logger.Stderr.Write([]byte(err.Error()))
		return false
	}

	bfo := bufio.NewWriter(fo)
	
	err = tmpl.Execute(bfo, input)

	if err != nil {
		logger.Stderr.Write([]byte(err.Error()))
		return false
	} else {
		bfo.Flush()
		return true
	}
}
