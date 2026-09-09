package main

import (
	"fmt"
	"log"
	"ranga/internal/board"
)

var version string = ""

func main() {

	b := board.NewBoard()

	if err := b.ParseFEN("6r1/8/k7/8/3P4/P4K1R/8/3N4 b - - 0 1"); err != nil {
		log.Fatal(err)
	}

	b.Print()

	fmt.Println(b.FEN())

	// engine := uci.NewEngine(os.Stdin, os.Stdout, version)
	// engine.Run()

}
