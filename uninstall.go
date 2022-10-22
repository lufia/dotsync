package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func runUninstall(r *Repository, args []string, w io.Writer) error {
	f := NewFlagSet("uninstall", "[path ...]")
	force := f.Bool("f", false, "discard locally changes and remove")

	if err := f.Parse(args); err != nil {
		return err
	}
	targets := make(map[string]struct{})
	for _, arg := range f.Args() {
		s, err := filepath.Abs(arg)
		if err != nil {
			return err
		}
		targets[s] = struct{}{}
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
		if _, ok := targets[state.Target]; !ok {
			return nil
		}

		if !*force {
			h, err := ReadHash(state.Target)
			if err != nil {
				return err
			}
			ok, _, err := isModeEqual(state.Target, state.Mode)
			if err != nil {
				return err
			}
			if !ok || h != state.Hash {
				return fmt.Errorf("%s: locally modified; will not remove", state.Target)
			}
		}
		if err := remove(p); err != nil {
			return err
		}
		if err := remove(state.Target); err != nil {
			return err
		}
		return nil
	})
}

func remove(name string) error {
	if err := os.Remove(name); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}
