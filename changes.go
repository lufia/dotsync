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
		if h != state.Hash {
			fmt.Fprintln(w, state.Target)
		}
		return nil
	})
}
