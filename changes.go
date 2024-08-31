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
		ok, _, err := isModeEqual(state.Target, state.Mode)
		if err != nil {
			return err
		}
		if !ok || h != state.Hash {
			fmt.Fprintln(w, state.Target)
		}
		return nil
	})
}

func isModeEqual(file string, mode os.FileMode) (bool, os.FileMode, error) {
	fi, err := os.Stat(file)
	if err != nil {
		return false, 0, err
	}
	m := fi.Mode() & os.ModePerm
	return m == mode, m, nil
}
