package main

import (
	"flag"
	"fmt"
)

// XDG_STATE_HOME

const rootFlagKey = "test.r"

func NewFlagSet(name string) *flag.FlagSet {
	f := flag.NewFlagSet(name, flag.ExitOnError)
	f.Usage = func() {
		fmt.Fprintf(f.Output(), "usage: %s [options]\n", name)
		f.PrintDefaults()
	}
	if flag.Lookup("test.v") != nil {
		f.String(rootFlagKey, "", "root `dir`ectory")
	}
	return f
}

func RootDir(f *flag.FlagSet) string {
	root := f.Lookup(rootFlagKey)
	if root != nil && root.Value.String() != "" {
		return root.Value.String()
	}
	return defaultRootDir()
}

func defaultRootDir() string {
	return "."
}
