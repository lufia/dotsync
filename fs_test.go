package main

import (
	"os"
	"path/filepath"
	"strconv"
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
		attr, err := parseFileAttr(f.Name)
		if err != nil {
			t.Fatal(err)
		}
		name := filepath.Join(rootDir, attr.Name)
		s := os.ExpandEnv(string(f.Data))
		if err := writeFile(name, []byte(s), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(name, attr.Mode); err != nil {
			t.Fatal(err)
		}
	}
	return rootDir
}

type fileAttr struct {
	Name string
	Mode os.FileMode
}

func parseFileAttr(s string) (*fileAttr, error) {
	attr := fileAttr{
		Name: s,
		Mode: 0644,
	}
	if i := strings.IndexByte(s, '!'); i >= 0 {
		attr.Name = s[:i]
		m, err := strconv.ParseInt(s[i+1:], 8, 0)
		if err != nil {
			return nil, err
		}
		attr.Mode = os.FileMode(m)
	}
	return &attr, nil
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
	s1, err := os.Stat(golden)
	if err != nil {
		t.Fatal(err)
	}
	s2, err := os.Stat(actual)
	if err != nil {
		t.Fatal(err)
	}
	if s1.Mode() != s2.Mode() {
		t.Errorf("mismath want: mode=%o, got: mode=%o", s1.Mode(), s2.Mode())
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
