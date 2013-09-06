package main

import (
	"testing"
)

func equals(t *testing.T, desc string, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("%s: wantet %q, got %q", desc, a, b)
	}
}
