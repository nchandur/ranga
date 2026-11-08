package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"ranga/internal/uci"
)

var version string = ""

func main() {

	if len(os.Args) > 1 && os.Args[1] == "bench" {
		engine := uci.NewEngine(os.Stdin, os.Stdout, version)
		uci.Bench(engine, os.Stdout)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "genfile" {
		genFile(os.Args[2:])
		return
	}

	in := io.Reader(os.Stdin)
	if len(os.Args) > 1 {
		in = strings.NewReader(strings.Join(os.Args[1:], "\n") + "\nquit\n")
	}

	engine := uci.NewEngine(in, os.Stdout, version)
	engine.Run()
}

// handles the genfile mode
func genFile(args []string) {
	fs := flag.NewFlagSet("genfile", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	out := fs.String("out", "", "path of the file to write (required)")
	count := fs.Int("count", -1, "number of positions to generate (required)")
	seed := fs.Uint64("seed", 0, "random seed")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: ranga genfile -out <path> -count <n> [-seed <s>]")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
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

	n, err := engine.GenFENsToFile(*out, *count, *seed)
	if err != nil {
		fmt.Fprintln(os.Stderr, "genfile:", err)
		os.Exit(1)
	}
	fmt.Printf("wrote %d positions to %s\n", n, *out)
}
