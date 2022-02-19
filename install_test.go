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
		{
			// set permission
			file: "testdata/install/init.txtar",
			src:  "~/dotfiles/bin/ct",
			dest: "~/out/bin/ct",
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			dir, r := initFS(t, tt.file)
			args := []string{
				expandTilde(dir, tt.src),
				expandTilde(dir, tt.dest),
			}
			if err := runInstall(r, args, os.Stdout); err != nil {
				t.Fatal(err)
			}
			testFileContent(t, args[0], args[1])

			s, err := filepath.Rel(r.rootDir, args[0])
			if err != nil {
				t.Fatal(err)
			}
			file := filepath.Join(r.StateDir, "store", s)
			testFileContent(t, file+".golden", file)
		})
	}
}

func TestRunInstallDir(t *testing.T) {
	tests := []struct {
		file string
		src  string
		dest string
	}{
		{
			file: "testdata/install/dir.txtar",
			src:  "~/dotfiles/.exrc",
			dest: "~/out/",
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			dir, r := initFS(t, tt.file)
			args := []string{
				expandTilde(dir, tt.src),
				expandTilde(dir, tt.dest),
			}
			if err := runInstall(r, args, os.Stdout); err != nil {
				t.Fatal(err)
			}
			target := filepath.Join(args[1], filepath.Base(args[0]))
			testFileContent(t, args[0], target)

			s, err := filepath.Rel(r.rootDir, args[0])
			if err != nil {
				t.Fatal(err)
			}
			file := filepath.Join(r.StateDir, "store", s)
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
			dir, r := initFS(t, tt.file)
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
