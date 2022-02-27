package main

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/txtar"
)

func initFS(t testing.TB, file string) (string, *Repository) {
	t.Helper()
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	return initFSR(t, file, data)
}

func initFSR(t testing.TB, file string, data []byte) (string, *Repository) {
	a := txtar.Parse(data)
	rootDir := t.TempDir()
	os.Setenv("TEST_DIR", rootDir)
	t.Cleanup(func() {
		os.Unsetenv("TEST_DIR")
	})
	for _, f := range a.Files {
		attr, err := parseFileAttr(f.Name)
		if err != nil {
			t.Fatalf("%s: cannot parse '%s'", file, f.Name)
		}
		name := filepath.Join(rootDir, attr.Name)
		if attr.Mode.IsDir() {
			mkdirFatal(t, name, attr.Mode)
			continue
		}
		s := os.ExpandEnv(string(f.Data))
		opts := FileOptions{
			MkdirAll: true,
		}
		if err := writeFile(name, []byte(s), opts); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(name, attr.Mode); err != nil {
			t.Fatal(err)
		}
	}
	return rootDir, &Repository{
		StateDir: filepath.Join(rootDir, ".local/state/dotsync"),
		rootDir:  filepath.Join(rootDir, "dotfiles"),
	}
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
		s = s[i+1:]
		if strings.HasPrefix(s, "d") {
			s = s[1:]
			attr.Mode = os.ModeDir | 0755
		}
		m, err := strconv.ParseInt(s, 8, 0)
		if err != nil {
			return nil, err
		}
		attr.Mode &^= os.ModePerm
		attr.Mode |= os.FileMode(m)
	}
	return &attr, nil
}

func testFileContent(t testing.TB, golden, actual string) {
	t.Helper()
	wantBytes, err := os.ReadFile(golden)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			_, err := os.ReadFile(actual)
			if !errors.Is(err, os.ErrNotExist) {
				t.Errorf("file '%s' should not exist", actual)
			}
			return
		}
		t.Fatalf("read %s: %v", golden, err)
	}
	want := string(wantBytes)
	got := readFileFatal(t, actual)
	a := strings.Split(want, "\n")
	b := strings.Split(got, "\n")
	if diff := cmp.Diff(a, b); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
	s1 := statFatal(t, golden)
	s2 := statFatal(t, actual)
	if s1.Mode() != s2.Mode() {
		t.Errorf("mismath want: mode=%o, got: mode=%o", s1.Mode(), s2.Mode())
	}
}

func readFileFatal(t testing.TB, name string) string {
	t.Helper()
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("read %v: %v", name, err)
	}
	return string(data)
}

func statFatal(t testing.TB, name string) os.FileInfo {
	t.Helper()
	s, err := os.Stat(name)
	if err != nil {
		t.Fatalf("stat %s: %v", name, err)
	}
	return s
}

func mkdirFatal(t testing.TB, name string, perm os.FileMode) {
	t.Helper()
	err := os.MkdirAll(name, perm)
	if err != nil {
		t.Fatalf("mkdir %s: %v", name, err)
	}
}

func expandTilde(dir, s string) string {
	if s == "~" {
		return dir
	}
	s = filepath.ToSlash(s)
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
