package hce

import "ranga/internal/board"

type HCE struct{}

func (e *HCE) Evaluate(b *board.Board) int {
	score := 0

	for _, piece := range b.Mailbox {
		score += board.PieceValue[piece]
	}

	return score
}
