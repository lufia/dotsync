package main

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestRunPull(t *testing.T) {
	tests := []struct {
		file  string
		slugs []string
	}{
		{
			file:  "testdata/pull/updated.txtar",
			slugs: []string{".exrc", "lib/profile"},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			dir := initFS(t, tt.file)
			stateDir := filepath.Join(dir, ".local/state/dotsync")
			r := &Repository{
				StateDir: stateDir,
				rootDir:  filepath.Join(dir, "dotfiles"),
			}
			if err := runPull(r, nil, os.Stdout); err != nil {
				t.Fatal(err)
			}
			for _, slug := range tt.slugs {
				file := filepath.Join(dir, "out", slug)
				testFileContent(t, file+".golden", file)
			}
		})
	}
}
