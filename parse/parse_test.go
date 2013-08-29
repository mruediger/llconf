package parse

import (
	"testing"
	"fmt"
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

func TestPromiseFile(t *testing.T) {
	file,_  := os.Open( "./setup.cp" )
	bufin := bufio.NewReader( file )


	promises,_ := ReadPromises( bufin )
	len_wanted := 16
	
	if len( promises ) != len_wanted {
		t.Errorf("missing promisses: want %d, got %d", len_wanted, len(promises))
	}
}


func TestParser(t *testing.T) {
	file,e := os.Open( "./setup.cp" )
	if e != nil {
		panic(e)
	}
	bufin := bufio.NewReader( file )
	promises,_ := ParsePromises( bufin )

	for k,v := range(promises) {
		if k == "done" {
			fmt.Printf("%s: %s\n", k,v)
		}
	}
}
