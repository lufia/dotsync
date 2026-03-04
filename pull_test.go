package main

import (
	"os"
	"testing"
)

func TestRunPull(t *testing.T) {
	tests := []struct {
		script string
		label  string
	}{
		{"testdata/pull/arrived.script", "copy updated files"},
		{"testdata/pull/missing.script", "missing files show warnings but doesn't stop"},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			testRunFunc(t, tt.script, os.Stdout)
		})
	}
}
