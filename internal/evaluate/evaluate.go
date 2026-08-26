package evaluate

import (
	"ranga/internal/board"
	"ranga/internal/nnue"
)

type Evaluator struct{}

func (e *Evaluator) Evaluate(b *board.Board) int {
	score := 0

	for _, piece := range b.Mailbox {
		score += board.PieceValue[piece]
	}

	return score
}

type NNUEEvaluator struct{}

func (e *NNUEEvaluator) Evaluate(b *board.Board) int {
	return nnue.Evaluate(b)
}
