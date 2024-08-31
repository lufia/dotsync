package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func runInstall(r *Repository, args []string, w io.Writer) error {
	f := NewFlagSet("install", "dotfile [...] path")
	force := f.Bool("f", false, "indicate to replace old file")

	if err := f.Parse(args); err != nil {
		return err
	}
	args = f.Args()
	if len(args) < 2 {
		f.Usage()
		os.Exit(2)
	}
	to := args[len(args)-1]
	args = args[:len(args)-1]
	if len(args) >= 2 {
		ok, err := isDir(to)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("%s not a directory", to)
		}
	}
	for _, f := range args {
		err := r.CopyFile(to, f, CopyFileOptions{
			Overwrite: *force,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) CopyFile(dest, p string, opts CopyFileOptions) error {
	slug, err := r.Slug(p)
	if err != nil {
		return err
	}

	dest, err = filepath.Abs(dest)
	if err != nil {
		return err
	}
	ok, err := isDir(dest)
	if err != nil {
		return err
	}
	if ok {
		dest = filepath.Join(dest, filepath.Base(p))
	}
	h, mode, err := CopyFile(dest, p, opts)
	if err != nil {
		return err
	}

	s := fmt.Sprintf("%x %o %s\n", h, mode, dest)
	file := r.StateFile(slug)
	return writeFile(file, []byte(s), FileOptions{
		MkdirAll: true,
	})
}

func isDir(name string) (bool, error) {
	fi, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return fi.Mode().IsDir(), nil
}

type FileOptions struct {
	MkdirAll bool
	Perm     os.FileMode
}

func writeFile(name string, data []byte, opts FileOptions) error {
	dir := filepath.Dir(name)
	if opts.MkdirAll {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	perm := os.FileMode(0644)
	if opts.Perm != 0 {
		perm = opts.Perm
	}
	return os.WriteFile(name, data, perm)
}
