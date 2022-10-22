package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func runPush(r *Repository, args []string, w io.Writer) error {
	f := NewFlagSet("push", "[path ...]")
	dryRun := f.Bool("n", false, "dry run")
	_ = dryRun

	if err := f.Parse(args); err != nil {
		return err
	}
	dir := filepath.Join(r.StateDir, "store")
	err := filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
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
			fmt.Printf("cp %q %q\n", state.Target, state.Source)
		}
		ok, mode, err := isModeEqual(state.Target, state.Mode)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Printf("chmod %o %q\n", mode, state.Source)
		}
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}
