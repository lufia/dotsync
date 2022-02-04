package main

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestRunInstall(t *testing.T) {
	tests := []struct {
		file string
		src  string
		dest string
	}{
		{
			file: "testdata/install/init.txtar",
			src:  "~/dotfiles/.exrc",
			dest: "~/out/.exrc",
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			dir := initFS(t, tt.file)
			stateDir := filepath.Join(dir, ".local/state")
			r := &Repository{
				StateDir: stateDir,
				rootDir:  filepath.Join(dir, "dotfiles"),
			}

			args := []string{
				expandTilde(dir, tt.src),
				expandTilde(dir, tt.dest),
			}
			if err := runInstall(r, args, os.Stdout); err != nil {
				t.Fatal(err)
			}
			testFileContent(t, args[0], args[1])

			file := filepath.Join(stateDir, "dotsync/store/.exrc")
			testFileContent(t, file+".golden", file)
		})
	}
}

func TestRunInstallErr(t *testing.T) {
	tests := []struct {
		file string
		src  string
		dest string
		err  error
	}{
		{
			file: "testdata/install/exist.txtar",
			src:  "~/dotfiles/.exrc",
			dest: "~/out/.exrc",
			err:  os.ErrExist,
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			dir := initFS(t, tt.file)
			stateDir := filepath.Join(dir, ".local/state")
			r := &Repository{
				StateDir: stateDir,
				rootDir:  filepath.Join(dir, "dotfiles"),
			}

			args := []string{
				expandTilde(dir, tt.src),
				expandTilde(dir, tt.dest),
			}
			err := runInstall(r, args, os.Stdout)
			if !errors.Is(err, tt.err) {
				t.Errorf("runInstall(%v): err = %v; want %v", args, err, tt.err)
			}
			testFileContent(t, args[1]+".golden", args[1])
		})
	}
}
