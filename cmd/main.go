package main

import (
	"fmt"
	"ranga/internal/board"
)

var version string = ""

func main() {

	b := board.NewBoard()
	b.ParseFEN("Nnb3nr/r1pp4/3k1pp1/pPb3p1/3p1P2/5Q1P/P2BP3/2R1KBNR b K - 0 1")

	fmt.Printf("%b\n", b.Castle)

	// engine := uci.NewEngine(os.Stdin, os.Stdout, version)
	// engine.Run()

}
