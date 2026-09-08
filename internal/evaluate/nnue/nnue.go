package nnue

import "ranga/internal/board"

type NNUE struct {
	Network
}

func (n NNUE) Evaluate(_ *board.Board) int {
	return 0
}
