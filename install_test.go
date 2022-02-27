package main

import (
	"os"
	"testing"
)

func TestRunInstall(t *testing.T) {
	tests := []struct {
		script string
		label  string
	}{
		{"testdata/install/regular.script", "copy a regular file"},
		{"testdata/install/perm.script", "copy an executable file with its permission"},
		{"testdata/install/dir.script", "copy a file into a directory"},
		{"testdata/install/exist.script", "occurs an error when target is exist"},
		{"testdata/install/nodir.script", "occurs an error when parent directory is not exist"},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			testRunFunc(t, tt.script, os.Stdout)
		})
	}
}
