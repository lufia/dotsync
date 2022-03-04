package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

var commands = map[string]func(r *Repository, args []string, w io.Writer) error{
	"repo":      runRepo,
	"install":   runInstall,
	"changes":   runChanges,
	"pull":      runPull,
	"uninstall": runUninstall,
	"export":    runExport,
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("dotsync: ")
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(2)
	}
	f, ok := commands[args[0]]
	if !ok {
		usage()
		os.Exit(2)
	}
	r, err := OpenRepository()
	if err != nil {
		log.Fatal(err)
	}
	if err := f(r, args[1:], os.Stdout); err != nil {
		if e, ok := err.(interface {
			HasAlreadyShown() bool
		}); !ok || !e.HasAlreadyShown() {
			log.Fatal(err)
		}
		os.Exit(1)
	}
}

func usage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, "usage: %s [options] [commands]\n", os.Args[0])
	flag.PrintDefaults()

	fmt.Fprintln(w, "\navailable commands:\n")
	var cmds []string
	for s := range commands {
		cmds = append(cmds, s)
	}
	sort.Strings(cmds)
	for _, s := range cmds {
		fmt.Fprintf(w, "  %s\n", s)
	}
}
