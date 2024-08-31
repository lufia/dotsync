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
			h, mode, err := CopyFile(state.Source, state.Target, CopyFileOptions{
				Overwrite: true,
			})
			if err != nil {
				return err
			}
			state.Hash = string(h)
			state.Mode = mode
		}
		ok, mode, err := isModeEqual(state.Target, state.Mode)
		if err != nil {
			return err
		}
		if !ok {
			if err := os.Chmod(state.Source, mode); err != nil {
				return err
			}
			state.Mode = mode
		}
		s := fmt.Sprintf("%s %o %s\n", state.Hash, state.Mode, state.Target)
		return writeFile(p, []byte(s), FileOptions{})
	})
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}
