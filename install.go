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
		return err
	}
	args = f.Args()
	if len(args) != 2 {
		f.Usage()
		os.Exit(2)
	}
	return r.CopyFile(args[1], args[0], false)
}

func (r *Repository) CopyFile(dest, p string, overwrite bool) error {
	slug, err := r.Slug(p)
	if err != nil {
		return err
	}
	dest, err = filepath.Abs(dest)
	if err != nil {
		return err
	}
	fin, err := os.Open(p)
	if err != nil {
		return err
	}
	defer fin.Close()

	dir := filepath.Dir(dest)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	flags := os.O_WRONLY | os.O_CREATE
	if overwrite {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}
	fout, err := os.OpenFile(dest, flags, 0644)
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
	s := fmt.Sprintf("%x %s\n", h.Sum(nil), dest)
	file := r.StateFile(slug)
	return writeFile(file, []byte(s), 0644)
}

func writeFile(name string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(name, data, perm)
}
