package promise

import (
	"bytes"
	"testing"
)


func TestRestartPromise(t *testing.T) {
	promise := RestartPromise{}
	var sout, serr bytes.Buffer
	cmd, _ := promise.restartLLConf("/bin/echo", []string{"-n","hello world"}, &sout, &serr)
	cmd.Wait()
	equals(t, "hello world", sout.String())
}
