package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/txtar"
)

func initFS(t testing.TB, file string) string {
	t.Helper()
	a, err := txtar.ParseFile(file)
	if err != nil {
		t.Fatal(err)
	}
	rootDir := t.TempDir()
	for _, f := range a.Files {
		dir := filepath.Join(rootDir, filepath.Dir(f.Name))
		if err := os.MkdirAll(dir, 0755); err != nil {
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

func expandTilde(dir, s string) string {
	if s == "~" {
		return dir
	}
	if strings.HasPrefix(s, "~/") {
		return filepath.Join(dir, s[2:])
	}
	return s
}
