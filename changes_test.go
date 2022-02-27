package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRunChanges(t *testing.T) {
	tests := []struct {
		script string
		label  string
		want   string
	}{
		{
			script: "testdata/changes/clean.script",
			label:  "no output if files is not modified locally",
			want:   "",
		},
		{
			script: "testdata/changes/dirty.script",
			label:  "prints modified filenames",
			want: strings.Join([]string{
				"~/out/.exrc",
				"~/out/lib/newstime",
				"~/out/lib/profile",
			}, "\n") + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			var w bytes.Buffer
			dir := testRunFunc(t, tt.script, &w)
			a := expandTildeSlice(dir, strings.Split(tt.want, "\n"))
			b := strings.Split(w.String(), "\n")
			if diff := cmp.Diff(a, b); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
