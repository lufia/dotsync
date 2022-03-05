package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func runExport(r *Repository, args []string, w io.Writer) error {
	f := NewFlagSet("export", "")

	if err := f.Parse(args); err != nil {
		return err
	}
	dir := JoinName(r.StateDir, "store")
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
		fmt.Fprintf(w, "install %q %q\n", state.Source, state.Target)
		return nil
	})
}
