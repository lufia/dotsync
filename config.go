package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

func userStateDir() (string, error) {
	var dir string

	switch runtime.GOOS {
	case "plan9":
		dir = os.Getenv("home")
		if dir == "" {
			return "", errors.New("$home is not defined")
		}
		dir = filepath.Join(dir, "lib/state")
	default:
		dir = os.Getenv("XDG_STATE_HOME")
		if dir == "" {
			dir = os.Getenv("HOME")
			if dir == "" {
				return "", errors.New("neither $XDG_STATE_HOME nor $HOME are defined")
			}
			dir = filepath.Join(dir, ".local/state")
		}
	}
	return dir, nil
}

type Repository struct {
	StateDir string
	rootDir  string
}

var ErrNotInitialized = errors.New("repository is not initialized")

func OpenRepository() (*Repository, error) {
	dir, err := userStateDir()
	if err != nil {
		return nil, err
	}
	return &Repository{
		StateDir: dir,
	}, nil
}

func (r *Repository) Path() (string, error) {
	if r.rootDir != "" {
		return r.rootDir, nil
	}

	file := filepath.Join(r.StateDir, "dotsync", "repo")
	b, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrNotInitialized
		}
		return "", err
	}
	s := string(bytes.TrimSpace(b))
	if s == "" {
		return "", ErrNotInitialized
	}
	r.rootDir = s
	return r.rootDir, nil
}

func (r *Repository) PutPath(p string) error {
	p, err := filepath.Abs(p)
	if err != nil {
		return err
	}
	file := filepath.Join(r.StateDir, "dotsync", "repo")
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(file, []byte(p+"\n"), 0644); err != nil {
		return err
	}
	r.rootDir = p
	return nil
}

func (r *Repository) Slug(p string) (string, error) {
	p, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	dir, err := r.Path()
	if err != nil {
		return "", err
	}
	return filepath.Rel(dir, p)
}

func (r *Repository) StateFile(slug string) string {
	return filepath.Join(r.StateDir, "dotsync", "store", slug)
}
