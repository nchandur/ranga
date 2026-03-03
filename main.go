package main

import (
	"fmt"

	"github.com/nchandur/ranga/board"
)

func init() {
	board.Init120to64()
	board.InitBitMasks()
	board.InitHashKeys()
}

func main() {
	start := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	b := board.NewBoard()

	b.ParseFEN(start)

	fmt.Printf("%x\n", b.PositionKey)

}
