package main

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestRunUninstall(t *testing.T) {
	tests := []struct {
		file string
		args []string
	}{
		{
			file: "testdata/uninstall/init.txtar",
			args: []string{"~/out/.exrc"},
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
			args := expandTildeSlice(dir, tt.args)
			if err := runUninstall(r, args, os.Stdout); err != nil {
				t.Fatal(err)
			}

			outDir := filepath.Join(dir, "out")
			for _, arg := range args {
				slug, err := filepath.Rel(outDir, arg)
				if err != nil {
					t.Fatal(err)
				}
				testFileRemoved(t, arg)
				file := filepath.Join(stateDir, "store", slug)
				testFileRemoved(t, file)
			}
		})
	}
}

func testFileRemoved(t *testing.T, name string) {
	_, err := os.Stat(name)
	if err == nil {
		t.Errorf("%s should be removed; but it is still exist", name)
		return
	}
	if !os.IsNotExist(err) {
		t.Fatalf("%s: %v", name, err)
	}
}
