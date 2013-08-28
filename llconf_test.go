package main

import (
	"testing"
	"strings"
)

func equals(t *testing.T, desc string, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("%s: wantet %q, got %q", desc, a, b)
	}
}

func TestParseArguments(t *testing.T) {
	args := []string{ "llconf", "serve", "hello", "world" }
	_,err := processArguments(args, nil)
	if err != nil {
		t.Errorf("%v\n", err)
	}
}

func TestProcessServeFlags(t *testing.T) {
	args := []string{ "hello", "world" }
	
	runescanner := strings.NewReader("(done)")
	cfg,err := processServeFlags("llconf", runescanner, args)
	if err != nil {
		t.Errorf("%v\n", err)
	} else {
		equals(t, "cfg.Goal", cfg.(ServeConfig).Goal, "done")
		if cfg.(ServeConfig).Input == nil {
			t.Errorf("cfg.Input is nil")
		}
	}
}
