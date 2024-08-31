package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func runDiff(r *Repository, args []string, w io.Writer) error {
	f := NewFlagSet("diff", "[path ...]")

	if err := f.Parse(args); err != nil {
		return err
	}
	where := make(map[string]struct{})
	for _, arg := range f.Args() {
		dest, err := filepath.Abs(arg)
		if err != nil {
			return err
		}
		where[dest] = struct{}{}
	}
	dir := filepath.Join(r.StateDir, "store")
	return filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		state, err := r.ReadState(p)
		if err != nil {
			return err
		}
		if _, ok := where[state.Target]; !ok && len(where) > 0 {
			return nil
		}
		h, err := ReadHash(state.Target)
		if err != nil {
			return err
		}
		ok, _, err := isModeEqual(state.Target, state.Mode)
		if err != nil {
			return err
		}
		if ok && h == state.Hash {
			return nil
		}
		return diffFiles(w, state.Source, state.Target)
	})
}

func diffFiles(w io.Writer, f1, f2 string) error {
	t1, err := os.ReadFile(f1)
	if err != nil {
		return err
	}
	t2, err := os.ReadFile(f2)
	if err != nil {
		return err
	}
	d := diffmatchpatch.New()
	diffs := d.DiffMain(string(t1), string(t2), true)
	fmt.Fprintln(w, d.DiffPrettyText(diffs))
	return nil
}
