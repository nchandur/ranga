package main

import (
	"github.com/nchandur/ranga/board"
)

func init() {
	board.Init120to64()
	board.InitBitMasks()
	board.InitHashKeys()
}

func main() {
}
