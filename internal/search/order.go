package search

import "ranga/internal/board"

// evaluates heuristic score to move to prioritize it during move ordering
func (s *Searcher) scoreMove(b *board.Board, move board.Move) int {

	if move.IsCapture() {
		var captured board.Piece

		if move.IsEnpass() {
			captured = board.WP
			if b.Side == board.White {
				captured = board.BP
			}
		} else {
			captured = b.Mailbox[move.Target()]
		}

		return MVVLVA[move.Piece()][captured] + 10000
	}

	if move.Promoted() == board.WQ || move.Promoted() == board.BQ {
		return 9000
	}
	return 8000

}

// orders the legal moves in given move list
func (s *Searcher) sortMove(b *board.Board, ml *board.MoveList) {

	score := make([]int, ml.Count)
	for i := 0; i < ml.Count; i++ {
		score[i] = s.scoreMove(b, ml.Moves[i])
	}

	for i := 1; i < ml.Count; i++ {
		keyMove := ml.Moves[i]
		keyScore := score[i]
		j := i - 1

		for j >= 0 && score[j] < keyScore {
			ml.Moves[j+1] = ml.Moves[j]
			score[j+1] = score[j]
			j--
		}
		ml.Moves[j+1] = keyMove
		score[j+1] = keyScore
	}
}
