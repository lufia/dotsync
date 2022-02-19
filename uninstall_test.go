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
			dir, r := initFS(t, tt.file)
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
				file := filepath.Join(r.StateDir, "store", slug)
				testFileRemoved(t, file)
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
