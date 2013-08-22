package llconf

import (
	"testing"
	"fmt"
	"bufio"
	"os"
)

func TestParser(t *testing.T) {
	file,_ := os.Open( "./setup.cp" )
	bufin := bufio.NewReader( file )
	promises := ParsePromises( bufin )

	for k,v := range(promises) {
		fmt.Printf("%s: %s\n", k,v)
	}
}
