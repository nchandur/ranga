package nnue

import (
	"ranga/internal/board"
)

type NNUE struct{}

func (e *NNUE) Evaluate(b *board.Board) int {
	return evaluate(b)
}
