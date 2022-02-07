package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func runChanges(r *Repository, args []string, w io.Writer) error {
	f := NewFlagSet("changes", "")

	if err := f.Parse(args); err != nil {
		return err
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
		h, err := ReadHash(state.Target)
		if err != nil {
			return err
		}
		modeDiff, err := isModeChanged(state.Target, state.Source)
		if err != nil {
			return err
		}
		if modeDiff || h != state.Hash {
			fmt.Fprintln(w, state.Target)
		}
		return nil
	})
}

func isModeChanged(s1, s2 string) (bool, error) {
	f1, err := os.Stat(s1)
	if err != nil {
		return false, err
	}
	f2, err := os.Stat(s2)
	if err != nil {
		return false, err
	}
	m1 := f1.Mode() & os.ModePerm
	m2 := f2.Mode() & os.ModePerm
	return m1 != m2, nil
}
