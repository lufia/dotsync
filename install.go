package main

import (
	"crypto/sha256"
	"errors"
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

type CopyFileOptions struct {
	MkdirAll  bool
	Overwrite bool
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
	fin, err := os.Open(p)
	if err != nil {
		return err
	}
	defer fin.Close()

	ok, err := isDir(dest)
	if err != nil {
		return err
	}
	if ok {
		dest = filepath.Join(dest, filepath.Base(p))
	}
	dir := filepath.Dir(dest)
	if opts.MkdirAll {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	flags := os.O_WRONLY | os.O_CREATE
	if opts.Overwrite {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}
	fi, err := os.Stat(p)
	if err != nil {
		return err
	}
	mode := fi.Mode() & os.ModePerm
	fout, err := os.OpenFile(dest, flags, mode)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("'%s' file does not exist", dest)
		}
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("'%s' file already exist", dest)
		}
		return err
	}
	defer fout.Close()

	h := sha256.New()
	o := io.MultiWriter(h, fout)
	io.Copy(o, fin)
	if err := fout.Sync(); err != nil {
		return err
	}
	if err := os.Chmod(dest, mode); err != nil {
		return err
	}
	s := fmt.Sprintf("%x %o %s\n", h.Sum(nil), mode, dest)
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
