package main

import (
	"testing"
)

func equals(t *testing.T, desc string, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("%s: wantet %q, got %q", desc, a, b)
	}
}

func TestProcessServeFlags(t *testing.T) {
	args := []string{ "-input-folder", "/tmp", "--incomming-folder","/var"}
	
	cfg,err := processServeFlags("llconf", args)

	if err != nil {
		t.Errorf("%v\n", err)
	} else {
		equals(t, "cfg.Goal", cfg.Goal, "done")
	}
}
