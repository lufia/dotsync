package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func runPull(r *Repository, args []string, w io.Writer) error {
	f := NewFlagSet("pull", "[path ...]")
	dryRun := f.Bool("n", false, "dry run")

	if err := f.Parse(args); err != nil {
		return err
	}
	_ = *dryRun
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
		h, err := ReadHash(state.Source)
		if err != nil {
			return err
		}
		if h == state.Hash {
			return nil
		}

		h, err = ReadHash(state.Target)
		if err != nil {
			return err
		}
		if h != state.Hash {
			log.Printf("%s: locally modified; will not overwrite\n", p)
			return nil
		}
		if *dryRun {
			fmt.Printf("cp %q %q\n", state.Source, state.Target)
			return nil
		}
		return r.CopyFile(state.Target, state.Source, true)
	})
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}
