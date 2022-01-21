package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func runRepo(args []string, w io.Writer) error {
	f := NewFlagSet("repo")
	dir := f.String("w", "", "repository `dir`ectory")

	if err := f.Parse(args); err != nil {
		return fmt.Errorf("can't parse flags: %w", err)
	}
	file := filepath.Join(RootDir(f), "repo")
	if *dir == "" {
		b, err := os.ReadFile(file)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("cannot read %s: %w", file, err)
		}
		b = bytes.TrimSpace(b)
		if len(b) > 0 {
			fmt.Fprintf(w, "%s\n", b)
		}
	} else {
		err := os.WriteFile(file, []byte(*dir+"\n"), 0644)
		if err != nil {
			return fmt.Errorf("cannot write %s: %w", file, err)
		}
	}
	return nil
}
