package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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
		StateDir: filepath.Join(dir, "dotsync"),
	}, nil
}

// Path returns the source repository path.
func (r *Repository) Path() (string, error) {
	if r.rootDir != "" {
		return r.rootDir, nil
	}

	file := filepath.Join(r.StateDir, "repo")
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
	file := filepath.Join(r.StateDir, "repo")
	if err := writeFile(file, []byte(p+"\n"), 0644); err != nil {
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
	return filepath.Join(r.StateDir, "store", slug)
}

type State struct {
	Hash   string
	Mode   os.FileMode
	Target string
	Source string
}

func (r *Repository) ReadState(file string) (*State, error) {
	p, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	storeDir := filepath.Join(r.StateDir, "store")
	slug, err := filepath.Rel(storeDir, p)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	s := strings.TrimSpace(string(data))
	a := strings.SplitN(s, " ", 3)
	if len(a) != 3 {
		return nil, fmt.Errorf("%s: state is corrupted", file)
	}
	dir, err := r.Path()
	if err != nil {
		return nil, err
	}
	mode, err := strconv.ParseInt(a[1], 8, 0)
	if err != nil {
		return nil, err
	}
	return &State{
		Hash:   a[0],
		Mode:   os.FileMode(mode),
		Target: a[2],
		Source: filepath.Join(dir, slug),
	}, nil
}

func ReadHash(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
