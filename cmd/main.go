package main

import (
	"log"
	"ranga/board"
)

var version string = ""

func main() {

	b := board.NewBoard()

	if err := b.ParseFEN(board.START); err != nil {
		log.Fatal(err)
	}

	b.Print()

}
