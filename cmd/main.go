package main

import (
	"os"
	"ranga/internal/uci"
)

var version string = ""

func main() {

	engine := uci.NewEngine(os.Stdin, os.Stdout, version)
	engine.Run()

}
