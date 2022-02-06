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
	os.Setenv("TEST_DIR", rootDir)
	t.Cleanup(func() {
		os.Unsetenv("TEST_DIR")
	})
	for _, f := range a.Files {
		name := filepath.Join(rootDir, f.Name)
		s := os.ExpandEnv(string(f.Data))
		if err := writeFile(name, []byte(s), 0644); err != nil {
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

func expandTildeSlice(dir string, paths []string) []string {
	a := make([]string, len(paths))
	for i, s := range paths {
		a[i] = expandTilde(dir, s)
	}
	return a
}
