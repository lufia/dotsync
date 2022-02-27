package main

import (
	"errors"
	"os"
	"strconv"
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
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			testRunFunc(t, tt.script, os.Stdout)
		})
	}
}

func TestRunInstallErr(t *testing.T) {
	tests := []struct {
		file string
		args []string
		err  error
	}{
		{
			file: "testdata/install/exist.txtar",
			args: []string{"~/dotfiles/.exrc", "~/out/.exrc"},
			err:  os.ErrExist,
		},
		{
			file: "testdata/install/nodir.txtar",
			args: []string{"~/dotfiles/.exrc", "~/out/dir/.exrc"},
			err:  os.ErrNotExist,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			dir, r := initFS(t, tt.file)
			args := expandTildeSlice(dir, tt.args)
			err := runInstall(r, args, os.Stdout)
			if !errors.Is(err, tt.err) {
				t.Errorf("runInstall(%v): err = %v; want %v", args, err, tt.err)
			}
			testFileContent(t, args[1]+".golden", args[1])
		})
	}
}
