package llconf

import (
	"testing"
	"strings"
	"reflect"
	"os"
	"bufio"
)

func TestReadPromise(t *testing.T) {
	var tests = []struct {
		input string
		want UnparsedPromise
	}{
		{ "(done)",
			UnparsedPromise{ "done", []UnparsedPromise{}, []Argument{}}},
		{ "commented promies (done)",
			UnparsedPromise{ "done", []UnparsedPromise{}, []Argument{}}},
		{ "verbose (commented) promise",
			UnparsedPromise{ "commented", []UnparsedPromise{}, []Argument{}}},
		{ "(and (bash) (vim))",
			UnparsedPromise{ "and", []UnparsedPromise{
				UnparsedPromise{ "bash", []UnparsedPromise{}, []Argument{}},
				UnparsedPromise{ "vim", []UnparsedPromise{}, []Argument{}},
			}, []Argument{}}},
	}
	for _,c := range tests {
		got := ReadPromises( strings.NewReader(c.input ) )[0]
		if ! reflect.DeepEqual(got, c.want) {
			t.Errorf("ReadPromises(%q) == %q, want %q", c.input, got, c.want)
		}
	}
}

func TestPromiseFile(t *testing.T) {
	file,_  := os.Open( "./setup.cp" )
	bufin := bufio.NewReader( file )


	promises := ReadPromises( bufin )
	len_wanted := 14
	
	if len( promises ) != len_wanted {
		t.Errorf("missing promisses: want %d, got %d", len_wanted, len(promises))
	}
}
