package main

import (
	"log"
	"ranga/internal/board"
	"ranga/internal/evaluate/nnue"
)

var version string = ""

func main() {

	b := board.NewBoard()
	b.ParseFEN(board.START)

	net, err := nnue.LoadNetwork("/Users/nchandur/workspace/ranga/data/beans.bin")

	if err != nil {
		log.Fatalf("loading network: %v", err)
	}
	nn := &nnue.NNUE{Network: *net}
	nn.Reset(&b)

	nn.PrintDebugInfo(&b, "starting")

	b.ParseFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	nn.Reset(&b)

	nn.PrintDebugInfo(&b, "kiwipete")

	// engine := uci.NewEngine(os.Stdin, os.Stdout, version)
	// engine.Run()

}
