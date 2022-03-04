package main

import (
	"os"
	"testing"
)

func TestRunUninstall(t *testing.T) {
	tests := []struct {
		script string
		label  string
		files  []string
	}{
		{
			script: "testdata/uninstall/single.script",
			label:  "remove a target file and its state",
			files: []string{
				"~/out/.exrc",
				"~/.local/state/dotsync/store/.exrc",
			},
		},
		{
			script: "testdata/uninstall/modified.script",
			label:  "occurs an error if file is modified",
		},
		{
			script: "testdata/uninstall/discard.script",
			label:  "remove a target file and its state if modified",
			files: []string{
				"~/out/.exrc",
				"~/.local/state/dotsync/store/.exrc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			dir := testRunFunc(t, tt.script, os.Stdout)
			removedFiles := expandTildeSlice(dir, tt.files)
			for _, arg := range removedFiles {
				testFileRemoved(t, arg)
			}
		})
	}
}

func testFileRemoved(t testing.TB, name string) {
	_, err := os.Stat(name)
	if err == nil {
		t.Errorf("%s should be removed; but it is still exist", name)
		return
	}
	if !os.IsNotExist(err) {
		t.Fatalf("stat %s: %v", name, err)
	}
}
