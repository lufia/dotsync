package main

import (
	"bytes"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRunChanges(t *testing.T) {
	tests := []struct {
		file string
		want string
	}{
		{
			file: "testdata/changes/init.txtar",
			want: "",
		},
		{
			file: "testdata/changes/modified.txtar",
			want: "~/out/.exrc\n~/out/lib/profile\n",
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			dir := initFS(t, tt.file)
			stateDir := filepath.Join(dir, ".local/state/dotsync")
			r := &Repository{
				StateDir: stateDir,
				rootDir:  filepath.Join(dir, "dotfiles"),
			}
			var w bytes.Buffer
			if err := runChanges(r, nil, &w); err != nil {
				t.Fatal(err)
			}
			a := expandTildeSlice(dir, strings.Split(tt.want, "\n"))
			b := expandTildeSlice(dir, strings.Split(w.String(), "\n"))
			if diff := cmp.Diff(a, b); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
