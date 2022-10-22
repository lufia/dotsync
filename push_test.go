package main

import (
	"os"
	"testing"
)

func TestRunPush(t *testing.T) {
	tests := []struct {
		script string
		label  string
	}{
		{"testdata/push/changed.script", "updates"},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			testRunFunc(t, tt.script, os.Stdout)
		})
	}
}
