package main

import (
	"errors"
	"fmt"
	"io"
)

func runRepo(r *Repository, args []string, w io.Writer) error {
	f := NewFlagSet("repo", "")
	dir := f.String("w", "", "repository `dir`ectory")

	if err := f.Parse(args); err != nil {
		return err
	}
	if *dir == "" {
		s, err := r.Path()
		if err != nil && !errors.Is(err, ErrNotInitialized) {
			return err
		}
		if s != "" {
			fmt.Fprintf(w, "%s\n", s)
		}
	} else {
		if err := r.PutPath(*dir); err != nil {
			return err
		}
	}
	return nil
}
