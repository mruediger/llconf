package parse

import (
	"testing"
	"bufio"
	"os"
	"strings"
	"reflect"
	"github.com/mruediger/llconf/promise"
)

func TestReadPromise(t *testing.T) {
	var tests = []struct {
		input string
		want UnparsedPromise
	}{
		{ "(done)",
			UnparsedPromise{ "done", []UnparsedPromise{}, []promise.Argument{}}},
		{ "commented promies (done)",
			UnparsedPromise{ "done", []UnparsedPromise{}, []promise.Argument{}}},
		{ "verbose (commented) promise",
			UnparsedPromise{ "commented", []UnparsedPromise{}, []promise.Argument{}}},
		{ "(and (bash) (vim))",
			UnparsedPromise{ "and", []UnparsedPromise{
				UnparsedPromise{ "bash", []UnparsedPromise{}, []promise.Argument{}},
				UnparsedPromise{ "vim", []UnparsedPromise{}, []promise.Argument{}},
			}, []promise.Argument{}}},

			
	}
	for _,c := range tests {
		promises,_ := ReadPromises( strings.NewReader(c.input ) )
		got := promises[0]
		if ! reflect.DeepEqual(got, c.want) {
			t.Errorf("ReadPromises(%q) == %q, want %q", c.input, got, c.want)
		}
	}
}

func TestReadArguments(t *testing.T) {
	var tests = []struct {
		input string
		want UnparsedPromise
	}{
		{ "(test \"echo\" [join \"hello\" \"world\"])",
			UnparsedPromise{ "test", []UnparsedPromise{}, []promise.Argument{ promise.Constant{"echo"},
				promise.JoinArgument{ []promise.Argument{ promise.Constant{ "hello" }, promise.Constant{ "world" }}}}}},
		{ "(test \"bla:fa:sel\")",
			UnparsedPromise{ "test", []UnparsedPromise{}, []promise.Argument{ promise.Constant{"bla:fa:sel"}}}},
		{ "(test [var:foo])",
			UnparsedPromise{ "test", []UnparsedPromise{}, []promise.Argument{ promise.VarGetter{"foo"}}}},
	}
	for _,test := range tests {
		promises,err := ReadPromises( strings.NewReader(test.input) )
		if err != nil {
			t.Errorf(err.Error())
		} else {
			got := promises[0]
			if ! reflect.DeepEqual(got, test.want) {
				t.Errorf("ReadPromises(%q) == %q, want %q", test.input, got, test.want)
			}
		}
	}
}

func TestReadJoiner(t *testing.T) {
	input := "join [arg:0] [env:test]]"
	reader :=  strings.NewReader(input ) 
	got,err := readArgument( reader, '[' )

	if err != nil {
		panic(err)
	}
	
	want := promise.JoinArgument{ []promise.Argument{ promise.ArgGetter{0}, promise.EnvGetter{"test"} }}
	if ! reflect.DeepEqual(got, want) {
		t.Errorf("readJoin(%q) == %q, wanted %q", input, got, want)
	}

	if reader.Len() > 0 {
		t.Errorf("%d chars left in reader",reader.Len())
	}
}

func TestPromiseFile(t *testing.T) {
	file,_  := os.Open( "./setup.cp" )
	bufin := bufio.NewReader( file )


	promises,_ := ReadPromises( bufin )
	len_wanted := 16
	
	if len( promises ) != len_wanted {
		t.Errorf("missing promisses: want %d, got %d", len_wanted, len(promises))
	}
}

