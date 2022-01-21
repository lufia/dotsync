package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/txtar"
)

func TestRunRepoRead(t *testing.T) {
	tests := []struct {
		file string
		s    string
	}{
		{file: "testdata/repo/empty.txtar", s: ""},
		{file: "testdata/repo/blank.txtar", s: ""},
		// TODO: valid
		// TODO: write
	}
	for _, tt := range tests {
		dir := initFS(t, tt.file)
		var w bytes.Buffer
		args := []string{"-test.r", dir}
		if err := runRepo(args, &w); err != nil {
			t.Fatal(err)
		}
		if s := w.String(); s != tt.s {
			t.Errorf("%s: got %s but want %s\n", tt.file, s, tt.s)
		}
	}
}

func initFS(t testing.TB, file string) string {
	a, err := txtar.ParseFile(file)
	if err != nil {
		t.Fatal(err)
	}
	rootDir := t.TempDir()
	for _, f := range a.Files {
		dir := filepath.Join(rootDir, filepath.Dir(f.Name))
		if err := os.MkdirAll(dir, 755); err != nil {
			t.Fatal(err)
		}
		name := filepath.Join(rootDir, f.Name)
		if err := os.WriteFile(name, f.Data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	return rootDir
}
