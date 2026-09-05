package nnue

import "ranga/internal/board"

type NNUE struct {
	*Network
	*Accumulator
}

func (n NNUE) Evaluate(b *board.Board) int {
	return n.Network.Evaluate(n.Accumulator, b.Side)
}
