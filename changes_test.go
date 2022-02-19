package main

import (
	"bytes"
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
			want: strings.Join([]string{
				"~/out/.exrc",
				"~/out/lib/newstime",
				"~/out/lib/profile",
			}, "\n") + "\n",
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			dir, r := initFS(t, tt.file)
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
