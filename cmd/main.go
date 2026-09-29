package main

import (
	"fmt"
	"ranga/internal/board"
)

var version string = ""

func main() {

	b := board.NewBoard()

	b.ParseFEN(board.START)
	// b.ParseFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")

	fmt.Println(b.Occupancies[2])

	// engine := uci.NewEngine(os.Stdin, os.Stdout, version)
	// engine.Run()

}
