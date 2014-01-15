package promise

import (
	"strings"
	"testing"
	"io"
	"os"
	"io/ioutil"
)


func runTest(template string,json string) string {
	fho, err := ioutil.TempFile("/tmp", "template-output")
	if err != nil { panic( err ) }
	fho.Close()
	output := fho.Name()

	fht, err := ioutil.TempFile("/tmp", "template-input")
	if err != nil { panic( err ) }
	io.WriteString(fht, template)
	fht.Close()
	template_file := fht.Name()

	promise := TemplatePromise{Constant(json), Constant(template_file), Constant(output)}

	ctx := NewContext()
	promise.Eval([]Constant{}, &ctx)


	bytes,err := ioutil.ReadFile(output)
	if err != nil { panic(err) }

	str := string(bytes)
	os.Remove(output)
	os.Remove(template_file)

	return str
}

func TestTemplatePromise(t *testing.T) {
	template := "{{ range . }} hallo {{.}} \n{{end}}"
	json := `["poolnode-01.ie01.d3sv.net","poolnode-02.ie01.d3sv.net"]`
	output := runTest(template,json)

	if ! strings.Contains(output, "poolnode") {
		t.Errorf("%q does not contain poolnode\n", output)
	}
}

func TestTemplatePromiseWithSingleVar(t *testing.T) {
	template := "hallo {{.}}\n"
	json := `'foobar'`
	output := runTest(template,json)

	if ! strings.Contains(output, "foobar") {
		t.Errorf("%q does not contain foobar\n", output)
	}

}
