package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// testRunFunc execute script and then checks all golden files that are contained in the script.
func testRunFunc(t testing.TB, script string, w io.Writer) string {
	t.Helper()
	data, err := os.ReadFile(script)
	if err != nil {
		t.Fatal(err)
	}
	// FIXME: it's better to use bytes.Cut.
	defs := bytes.SplitN(data, []byte{'\n'}, 2)
	if len(defs) < 2 {
		t.Fatalf("%s: should describe command line args in the first line", script)
	}
	a := strings.Fields(string(defs[0]))
	if len(a) < 1 {
		t.Fatalf("%s: the command line args cannot be an empty", script)
	}
	run, ok := commands[a[0]]
	if !ok {
		t.Fatalf("%s: '%s' is not defined in commands", script, a[0])
	}

	dir, r := initFSR(t, script, defs[1])
	args := expandTildeSlice(dir, a[1:])
	if err := run(r, args, w); err != nil {
		t.Fatalf("%s: %v", script, err)
	}

	err = filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(p) != ".golden" {
			return nil
		}
		name := p[:len(p)-7]
		testFileContent(t, name, p)
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", dir, err)
	}
	return dir
}
