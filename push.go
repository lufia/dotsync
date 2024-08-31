package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func runPush(r *Repository, args []string, w io.Writer) error {
	f := NewFlagSet("push", "[path ...]")
	dryRun := f.Bool("n", false, "dry run")

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
		if _, ok := where[state.Target]; !ok && len(where) > 0 {
			return nil
		}
		defer func() {
			delete(where, state.Target)
		}()
		h, err := ReadHash(state.Target)
		if err != nil {
			return err
		}
		if h != state.Hash {
			if *dryRun {
				fmt.Printf("cp %q %q\n", state.Target, state.Source)
				return nil
			}
			h, mode, err := CopyFile(state.Source, state.Target, CopyFileOptions{
				Overwrite: true,
			})
			if err != nil {
				return err
			}
			state.Hash = hex.EncodeToString(h)
			state.Mode = mode
		}
		ok, mode, err := isModeEqual(state.Target, state.Mode)
		if err != nil {
			return err
		}
		if !ok {
			if *dryRun {
				fmt.Printf("chmod %o %q\n", mode, state.Source)
				return nil
			}
			if err := os.Chmod(state.Source, mode); err != nil {
				return err
			}
			state.Mode = mode
		}
		if *dryRun {
			return nil
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
	if len(where) > 0 {
		errs := make([]error, len(where))
		for s := range where {
			errs = append(errs, fmt.Errorf("'%s' does not changed", s))
		}
		return errors.Join(errs...)
	}
	return nil
}
