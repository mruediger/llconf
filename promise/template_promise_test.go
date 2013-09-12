package promise

import (
	"strings"
	"testing"
	"io"
	"os"
	"io/ioutil"
)


func TestTemplatePromise(t *testing.T) {
	fho, err := ioutil.TempFile("/tmp", "template-output")
	if err != nil { panic( err ) }
	fho.Close()
	output := fho.Name()

	fht, err := ioutil.TempFile("/tmp", "template-input")
	if err != nil { panic( err ) }
	io.WriteString(fht, "{{ range . }} hallo {{.}} \n{{end}}")
	fht.Close()
	template := fht.Name()
	
	promise := TemplatePromise{ []Argument{
		Constant{`["poolnode-01.ie01.d3sv.net","poolnode-02.ie01.d3sv.net"]`},
		Constant{template},
		Constant{output}}}

	logger := Logger{Stdout:os.Stdout, Stderr: os.Stderr}
	promise.Eval([]Constant{}, &logger)


	bytes,err := ioutil.ReadFile(output)
	if err != nil { panic(err) }

	str := string(bytes)
	if ! strings.Contains(str, "poolnode") {
		t.Errorf("%q does not contain poolnode\n", str)
	}

	os.Remove(output)
	os.Remove(template)
}
