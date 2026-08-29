package hce

import "ranga/internal/board"

type HCE struct{}

func (h HCE) Evaluate(b *board.Board) int {
	score := 0

	for _, piece := range b.Mailbox {
		score += board.PieceValue[piece]
	}

	return score
}
