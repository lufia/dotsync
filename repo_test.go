package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/txtar"
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

func TestRunRepoWrite(t *testing.T) {
	tests := []struct {
		file string
		s    string
	}{
		{file: "testdata/repo/init.txtar", s: "src/dotfiles"},
		{file: "testdata/repo/replace.txtar", s: "src/dotfiles"},
	}
	for _, tt := range tests {
		dir := initFS(t, tt.file)
		var w bytes.Buffer
		args := []string{"-test.r", dir, "-w", "src/dotfiles"}
		if err := runRepo(args, &w); err != nil {
			t.Fatal(err)
		}
		if s := w.String(); s != "" {
			t.Errorf("%s: got %s but do not want any outputs\n", tt.file, s)
		}
		file := filepath.Join(dir, "repo")
		testFileContent(t, file+".golden", file)
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

func testFileContent(t testing.TB, golden, actual string) {
	t.Helper()
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(actual)
	if err != nil {
		t.Fatal(err)
	}
	a := strings.Split(string(want), "\n")
	b := strings.Split(string(got), "\n")
	if diff := cmp.Diff(a, b); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
