package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func runInstall(r *Repository, args []string, w io.Writer) error {
	f := NewFlagSet("install", "dotfile [path]")

	if err := f.Parse(args); err != nil {
		return fmt.Errorf("can't parse flags: %w", err)
	}
	args = f.Args()
	if len(args) != 2 {
		f.Usage()
		os.Exit(2)
	}
	slug, err := r.Slug(args[0])
	if err != nil {
		return err
	}
	fin, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer fin.Close()

	dir := filepath.Dir(args[1])
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	fout, err := os.OpenFile(args[1], os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer fout.Close()

	h := sha256.New()
	o := io.MultiWriter(h, fout)
	io.Copy(o, fin)
	if err := fout.Sync(); err != nil {
		return err
	}
	s := fmt.Sprintf("%x %s\n", h.Sum(nil), slug)
	file := r.StateFile(slug)
	return os.WriteFile(file, []byte(s), 0644)
}
