package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestRunRepoRead(t *testing.T) {
	tests := []struct {
		file string
		s    string
	}{
		{file: "testdata/repo/empty.txtar", s: ""},
		{file: "testdata/repo/blank.txtar", s: ""},
		{file: "testdata/repo/valid.txtar", s: "src/dotfiles\n"},
	}
	for _, tt := range tests {
		dir := initFS(t, tt.file)
		r := &Repository{
			StateDir: dir,
		}
		var args []string
		var w bytes.Buffer
		if err := runRepo(r, args, &w); err != nil {
			t.Fatal(err)
		}
		if s := w.String(); s != tt.s {
			t.Errorf("%s: got %s but want %s\n", tt.file, s, tt.s)
		}
	}
}

func TestRunRepoWrite(t *testing.T) {
	wdir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		file string
		s    string
		want string
	}{
		{
			file: "testdata/repo/empty.txtar",
			s:    "src/dotfiles",
			want: filepath.Join(wdir, "src", "dotfiles\n"),
		},
		{
			file: "testdata/repo/empty.txtar",
			s:    "/tmp/src/dotfiles",
			want: "/tmp/src/dotfiles\n",
		},
		{
			file: "testdata/repo/valid.txtar",
			s:    "src/dotfiles",
			want: filepath.Join(wdir, "src", "dotfiles\n"),
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			dir := initFS(t, tt.file)
			r := &Repository{
				StateDir: dir,
			}
			args := []string{"-w", tt.s}
			var w bytes.Buffer
			if err := runRepo(r, args, &w); err != nil {
				t.Fatal(err)
			}
			if s := w.String(); s != "" {
				t.Errorf("%s: got %s but do not want any outputs\n", tt.file, s)
			}
			file := filepath.Join(dir, "repo")
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			if s := string(data); s != tt.want {
				t.Errorf("runRepo(%v): repo = %q; want %q", args, s, tt.want)
			}
		})
	}
}
