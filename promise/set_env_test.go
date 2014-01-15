package promise

import (
	"os"
	"fmt"
	"bytes"
	"testing"
)


func TestSetEnvNew(t *testing.T) {
	old := SetEnv{}
	name := Constant("setenv")
	value := Constant("test")

	d,err := old.New([]Promise{DummyPromise{}}, []Argument{name,value})

	if err != nil {
		t.Errorf("(setenv) TestNew: %s", err.Error())
	} else {
		if d.(SetEnv).name != name {
			t.Errorf("(setenv) TestNew: env name not set")
		}
	}
}

func TestSetNewEval(t *testing.T) {
	arguments := []Argument{
		Constant("/bin/bash"),
		Constant("-c"),
		Constant("echo $setenv"),
	}
	exec := ExecPromise{ExecTest, arguments}

	var sout bytes.Buffer
	ctx := NewContext()
	ctx.Logger.Stdout = &sout
	ctx.Logger.Info = &sout

	name  := "setenv"
	value	:= "blafasel"
	s := SetEnv{Constant(name), Constant(value), exec}

	oldenv := fmt.Sprintf("%v",os.Environ())
	s.Eval([]Constant{}, &ctx)
	newenv := fmt.Sprintf("%v",os.Environ())

	if oldenv != newenv {
		t.Errorf("(setenv) changed overall environment")
	}

	if sout.String() != "/bin/bash -c echo $setenv\nblafasel\n" {
		t.Errorf("env name not present during execution")
	}
}
