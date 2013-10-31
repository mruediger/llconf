package promise

import (
	"time"
	"bytes"
	"testing"
)


func TestRestartPromise(t *testing.T) {
	promise := RestartPromise{}
	var sout, serr bytes.Buffer
	promise.restartLLConf("/bin/echo", []string{"-n","hello world"}, &sout, &serr)
	time.Sleep(time.Duration(1) * time.Second) // wait for the program to finish
	equals(t, "hello world", sout.String())
}
