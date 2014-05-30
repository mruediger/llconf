package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"testing"
	"time"
)

func TestWriteRunLog(t *testing.T) {
	starttime := time.Now().Local()
	endtime := time.Now().Local()
	path := "runlog.txt"
	writeRunLog(false, 1, 2, starttime, endtime, path)

	runlog, _ := ioutil.ReadFile(path)
	os.Remove(path)

	re := regexp.MustCompile("error, endtime=[0-9]+, duration=0.[0-9]+, c=1, t=2")
	if !re.Match(runlog) {
		t.Error("unexpected output")
	}
}
