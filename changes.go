package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func runChanges(r *Repository, args []string, w io.Writer) error {
	f := NewFlagSet("changes", "")

	if err := f.Parse(args); err != nil {
		return err
	}
	dir := filepath.Join(r.StateDir, "dotsync", "store")
	return filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		state, err := readState(p)
		if err != nil {
			return err
		}
		dir, err := r.Path()
		if err != nil {
			return err
		}
		target := filepath.Join(dir, state.Slug)
		f, err := os.Open(target)
		if err != nil {
			return err
		}
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return err
		}
		if s := fmt.Sprintf("%x", h.Sum(nil)); s != state.Hash {
			fmt.Fprintln(w, state.Slug)
		}
		return nil
	})
}

type State struct {
	Hash string
	Slug string
}

func readState(file string) (*State, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	s := strings.TrimSpace(string(data))
	a := strings.SplitN(s, " ", 2)
	if len(a) != 2 {
		return nil, fmt.Errorf("%s: state is corrupted", file)
	}
	return &State{
		Hash: a[0],
		Slug: a[1],
	}, nil
}
