package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"ranga/internal/uci"
)

var version string = ""

func main() {

	fs := flag.NewFlagSet("genfile", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	out := fs.String("out", "", "path of the file to write (required)")
	count := fs.Int("count", -1, "number of positions to generate (required)")
	seed := fs.Uint64("seed", 0, "random seed")
	book := fs.String("book", "", "opening book (EPD or FEN, one per line); omit to start from the standard position")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: ranga -out <path> -count <n> [-seed <s>] [-book <path>]")
		fs.PrintDefaults()
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		os.Exit(2)
	}

	if fs.NArg() > 0 || *out == "" || *count < 0 {
		fs.Usage()
		os.Exit(2)
	}

	engine := uci.NewEngine(os.Stdin, os.Stdout, version)

	n, err := engine.GenFENsToFile(*out, *count, *seed, *book)
	if err != nil {
		fmt.Fprintln(os.Stderr, "genfile:", err)
		os.Exit(1)
	}
	fmt.Printf("wrote %d positions to %s\n", n, *out)
}
