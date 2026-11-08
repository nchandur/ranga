package main

import (
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

	in := io.Reader(os.Stdin)
	if len(os.Args) > 1 {
		in = strings.NewReader(strings.Join(os.Args[1:], "\n") + "\nquit\n")
	}

	engine := uci.NewEngine(in, os.Stdout, version)
	engine.Run()
}
