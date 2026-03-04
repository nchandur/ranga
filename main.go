package main

import (
	"fmt"

	"github.com/nchandur/ranga/board"
)

func init() {
	board.InitBitMasks()
	board.InitHashKeys()
}

func main() {

	b := board.NewBoard()

	b.ParseFEN("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")
	b.GenPositionKey()

	fmt.Println(b.PositionKey)

}
