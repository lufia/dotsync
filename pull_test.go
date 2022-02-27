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
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			testRunFunc(t, tt.script, os.Stdout)
		})
	}
}
