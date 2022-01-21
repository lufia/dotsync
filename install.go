package main

import (
	"fmt"
	"io"
)

func runInstall(args []string, w io.Writer) error {
	f := NewFlagSet("install", "dotfile [path]")
	if err := f.Parse(args); err != nil {
		return fmt.Errorf("can't parse flags: %w", err)
	}
	return nil
}
