package main

import (
	"bytes"
	"testing"
)

func TestRunRepoRead(t *testing.T) {
	tests := []struct {
		dir string
		s   string
	}{
		{dir: "testdata/repo/empty", s: ""},
		// TODO: no file
		// TODO: valid
		// TODO: write
	}
	for _, tt := range tests {
		var w bytes.Buffer
		args := []string{"-test.r", tt.dir}
		if err := runRepo(args, &w); err != nil {
			t.Fatal(err)
		}
		if s := w.String(); s != tt.s {
			t.Errorf("%s: got %s but want %s\n", tt.dir, s, tt.s)
		}
	}
}
