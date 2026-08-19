package evaluate

import "ranga/internal/board"

type Evaluator struct{}

func (e *Evaluator) Evaluate(b *board.Board) int {
	score := 0

	for _, piece := range b.Mailbox {
		score += board.PieceValue[piece]
	}

	return score
}
