package main

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestRunRepoRead(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip()
	}

	tests := map[string]struct {
		file string
		s    string
	}{
		"repo file is not exist": {
			"testdata/repo/empty.txtar", "",
		},
		"repo file is blak file": {
			"testdata/repo/blank.txtar", "",
		},
		"repo file is already initialized": {
			"testdata/repo/valid.txtar", "src/dotfiles\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, r := initFS(t, tt.file)
			r.rootDir = ""
			var args []string
			var w bytes.Buffer
			if err := runRepo(r, args, &w); err != nil {
				t.Fatal(err)
			}
			if s := w.String(); s != tt.s {
				t.Errorf("%s: got %q but want %q\n", tt.file, s, tt.s)
			}
		})
	}
}

func TestRunRepoWrite(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip()
	}

	wdir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		file string
		s    string
		want string
	}{
		"repo file is not exist then writes relative path": {
			file: "testdata/repo/empty.txtar",
			s:    "src/dotfiles",
			want: filepath.Join(wdir, "src", "dotfiles\n"),
		},
		"repo file is not exist then writes absolute path": {
			file: "testdata/repo/empty.txtar",
			s:    "/tmp/src/dotfiles",
			want: "/tmp/src/dotfiles\n",
		},
		"repo file is already initialized": {
			file: "testdata/repo/valid.txtar",
			s:    "src/dotfiles",
			want: filepath.Join(wdir, "src", "dotfiles\n"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, r := initFS(t, tt.file)
			r.rootDir = ""
			args := []string{"-w", tt.s}
			var w bytes.Buffer
			if err := runRepo(r, args, &w); err != nil {
				t.Fatal(err)
			}
			if s := w.String(); s != "" {
				t.Errorf("%s: got %s but do not want any outputs\n", tt.file, s)
			}
			file := filepath.Join(r.StateDir, "repo")
			if s := readFileFatal(t, file); s != tt.want {
				t.Errorf("runRepo(%v): repo = %q; want %q", args, s, tt.want)
			}
		})
	}
}
