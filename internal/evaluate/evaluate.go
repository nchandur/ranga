package evaluate

import "ranga/internal/board"

type Evaluator interface {
	Evaluate(b *board.Board) int
}
