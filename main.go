package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var commands = map[string]func(args []string, w io.Writer) error{
	"repo":    runRepo,
	"install": runInstall,
}

func main() {
	log.SetFlags(0)
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
	if err := f(args[1:], os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func usage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, "usage: %s [options] [commands]\n", os.Args[0])
	flag.PrintDefaults()
}
