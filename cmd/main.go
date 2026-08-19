package main

import (
	"fmt"
	"log"
	"ranga/board"
)

var version string = ""

func main() {

	b := board.NewBoard()

	if err := b.ParseFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"); err != nil {
		log.Fatal(err)
	}

	b.Print()

	ml := board.NewMoveList()

	ml.GenerateMoves(&b)

	fmt.Println(ml)

}
