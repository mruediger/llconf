package promise

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	d := InDir{}

	if p,e := d.New([]Promise{DummyPromise{}},[]Argument{Constant("test")}); e != nil {
		t.Errorf("indir TestNew: %s", e.Error())
	} else {

		if p.(InDir).promise == nil {
			t.Errorf("indir.TestNew: promise is not set properly")
		}
	}
}

func TestEval(t *testing.T) {
	arguments := []Argument{
		Constant("/usr/bin/pwd"),
	}
	exec := ExecPromise{ExecTest, arguments}

	d := InDir{Constant("/var"),exec}

	var sout bytes.Buffer

	ctx := NewContext()
	ctx.Logger.Stdout = &sout
	ctx.Logger.Info = &sout

	if d.Eval([]Constant{}, &ctx) {
		if ctx.InDir == "/var" {
			t.Errorf("indir.TestEval: testdir creept outside scope")
		}
	} else {
		t.Errorf("indir.TestEval: eval didn't succeed")
	}

	out := sout.String()

	if out != "/usr/bin/pwd\n/var\n" {
		t.Errorf("exec not in right dir, found: %s", out)
	}
}
