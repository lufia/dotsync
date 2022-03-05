package main

import (
	"path/filepath"
)

func JoinName(elem ...string) string {
	s := filepath.Join(elem...)
	return filepath.ToSlash(s)
}
