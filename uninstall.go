package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func runUninstall(r *Repository, args []string, w io.Writer) error {
	f := NewFlagSet("uninstall", "path")
	force := f.Bool("f", false, "discard locally changes and remove")

	if err := f.Parse(args); err != nil {
		return err
	}
	args = f.Args()
	if len(args) != 1 {
		f.Usage()
		os.Exit(2)
	}
	target, err := filepath.Abs(args[0])
	if err != nil {
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
		if state.Target != target {
			return nil
		}

		h, err := ReadHash(state.Target)
		if err != nil {
			return err
		}
		ok, err := isModeEqual(state.Target, state.Mode)
		if err != nil {
			return err
		}
		if !*force && (!ok || h != state.Hash) {
			return fmt.Errorf("%s: locally modified; will not remove", state.Target)
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
