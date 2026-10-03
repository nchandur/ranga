package main

import (
	"os"
	"ranga/internal/uci"
)

var version string = ""

func main() {

	engine := uci.NewEngine(os.Stdin, os.Stdout, version)

	if len(os.Args) > 1 && os.Args[1] == "bench" {
		uci.Bench(engine, os.Stdout)
		return
	}

	engine.Run()

}
