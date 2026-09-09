package main

import (
	"fmt"
	"log"
	"ranga/internal/board"
	"ranga/internal/pgn"
)

var version string = ""

func main() {

	b := board.NewBoard()

	if err := b.ParseFEN("r3kbnr/1bqppppp/p1p5/8/4P3/1P1B4/P1P2PPP/RNBQK2R w KQkq - 1 8"); err != nil {
		log.Fatal(err)
	}

	b.Print()

	san := "O-O"

	move := pgn.ParseSAN(&b, san)

	fmt.Println("San: ", san, "-> Move: ", move.String())

	fmt.Println(move.IsCastle())
	fmt.Println(move.Source(), move.Target())

	// engine := uci.NewEngine(os.Stdin, os.Stdout, version)
	// engine.Run()

}
